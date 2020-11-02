package presentation

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	secretManager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	previousSecretManager "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	"github.com/gin-gonic/gin"
)

type EventHandler struct{}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

type EventCreateRequest struct {
	Type      string `json:"type"`
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
}

func (h EventHandler) Create(c *gin.Context) {
	projectNumber := os.Getenv("PROJECT_NUMBER")
	secretManagerClient, err := secretManager.NewClient(c)
	if err != nil {
		_ = c.Error(fmt.Errorf("secret manager client initialize errror: %w", err)).SetType(gin.ErrorTypePrivate)
		return
	}

	slackSigningSecretRequest := &previousSecretManager.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/slack-signing-secret/versions/latest", projectNumber),
	}
	slackSigningSecret, err := secretManagerClient.AccessSecretVersion(c, slackSigningSecretRequest)
	if err != nil {
		_ = c.Error(fmt.Errorf("fetch slack signing secret error: %w", err)).SetType(gin.ErrorTypePrivate)
		return
	}

	verifier, err := slack.NewSecretsVerifier(c.Request.Header, string(slackSigningSecret.Payload.Data))
	if err != nil {
		_ = c.Error(fmt.Errorf("slack secret verifier error: %w", err)).SetType(gin.ErrorTypePrivate)
		return
	}
	if err := verifier.Ensure(); err != nil {
		_ = c.Error(fmt.Errorf("ensure slack secret verifier error: %w", err)).SetType(gin.ErrorTypePrivate)
		return
	}

	bodyReader := io.TeeReader(c.Request.Body, &verifier)
	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		_ = c.Error(fmt.Errorf("read request body error: %w", err)).SetType(gin.ErrorTypePrivate)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
	if err != nil {
		_ = c.Error(fmt.Errorf("parse slack eventsAPI error: %w", err)).SetType(gin.ErrorTypePrivate)
		return
	}

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(body, &r)
		if err != nil {
			_ = c.Error(fmt.Errorf("slack url verification error: %w", err)).SetType(gin.ErrorTypePrivate)
			return
		}
		c.JSON(http.StatusOK, r.Challenge)
		return
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.ReactionAddedEvent:
			slackAccessTokenRequest := &previousSecretManager.AccessSecretVersionRequest{
				Name: fmt.Sprintf("projects/%s/secrets/slack-access-token/versions/latest", projectNumber),
			}
			slackAccessToken, err := secretManagerClient.AccessSecretVersion(c, slackAccessTokenRequest)
			if err != nil {
				_ = c.Error(fmt.Errorf("fetch slack signing secret error: %w", err)).SetType(gin.ErrorTypePrivate)
				return
			}

			api := slack.New(string(slackAccessToken.Payload.Data))
			conversationHistory, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{
				ChannelID: ev.Item.Channel,
				Inclusive: true,
				Latest:    ev.Item.Timestamp,
				Limit:     1,
			})
			if err != nil {
				_ = c.Error(fmt.Errorf("fetch conversation history error: %w", err)).SetType(gin.ErrorTypePrivate)
				return
			}
			logrus.Info(fmt.Sprintf("fetch message len %+v", len(conversationHistory.Messages)))
			for _, message := range conversationHistory.Messages {
				logrus.Info(fmt.Sprintf("fetch message zero %+v", message.Text))
			}

			buf := bytes.NewBuffer(nil)
			_ = gob.NewEncoder(buf).Encode(&ev)
			logrus.Info(fmt.Sprintf("ev: %+v", ev))

			bmsg, err := json.Marshal(ev)
			if err != nil {
				_ = c.Error(fmt.Errorf("fetch conversation history error: %w", err)).SetType(gin.ErrorTypePrivate)
				return
			}
			msg := &pubsub.Message{
				Data: bmsg,
			}

			client, err := pubsub.NewClient(c, "uchiyama-sandbox")
			if err != nil {
				_ = c.Error(fmt.Errorf("pubsub client error: %w", err)).SetType(gin.ErrorTypePrivate)
				return
			}

			topic := client.Topic("slack-event")
			if _, err := topic.Publish(c, msg).Get(c); err != nil {
				_ = c.Error(fmt.Errorf("could not publish message: %w", err)).SetType(gin.ErrorTypePrivate)
				return
			}
			c.JSON(http.StatusOK, nil)
		}
	default:
		_ = c.Error(
			errors.New(fmt.Sprintf("expected slack event not found, got %s", eventsAPIEvent.Type)),
		).SetType(gin.ErrorTypePublic)
		return
	}
}
