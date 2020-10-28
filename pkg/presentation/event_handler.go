package presentation

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"errors"
	"github.com/slack-go/slack"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EventHandler struct{}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

type EventCreateRequest struct {
	Type string `json:"type"`
	Token string `json:"token"`
	Challenge string `json:"challenge"`
}

func (h EventHandler) Create(c *gin.Context) {
	client, err := secretmanager.NewClient(c)
	if err != nil {
		jsonError(c, http.StatusInternalServerError, fmt.Errorf("secretmanager initialize errror: %w", err))
		return
	}
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/759555709793/secrets/slack-signing-secret/versions/latest",
	}
	result, err := client.AccessSecretVersion(c, req)
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}

	req = &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/759555709793/secrets/slack-access-token/versions/latest",
	}
	result2, err := client.AccessSecretVersion(c, req)
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}


	verifier, err := slack.NewSecretsVerifier(c.Request.Header, string(result.Payload.Data))
	if err != nil {
		jsonError(c, http.StatusInternalServerError, fmt.Errorf("slack signing secret error: %w",err))
		return
	}
	bodyReader := io.TeeReader(c.Request.Body, &verifier)
	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		jsonError(c, http.StatusInternalServerError, fmt.Errorf("read body error: %w",err))
		return
	}
	if err := verifier.Ensure(); err != nil {
		jsonError(c, http.StatusInternalServerError, fmt.Errorf("ensure secret error: %w",err))
		return
	}

	eventsAPIEvent, e := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
	if e != nil {
		jsonError(c, http.StatusInternalServerError, fmt.Errorf("initialize eventsAPIEvent error: %w", e))
		return
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(body, &r)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, r.Challenge)
		return
	} else if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.ReactionAddedEvent:
			api := slack.New(string(result2.Payload.Data))
			conversationHistory, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{
				ChannelID: ev.Item.Channel,
				Inclusive: true,
				Latest: ev.Item.Timestamp,
				Limit: 1,
			})
			if err != nil {
				jsonError(c, http.StatusInternalServerError, fmt.Errorf("fetch conversation history error: %w", err))
				return
			}
			logrus.Info(fmt.Sprintf("fetch message len %+v", len(conversationHistory.Messages)))
			for _, message := range conversationHistory.Messages {
				logrus.Info(fmt.Sprintf("fetch message zero %+v", message.Text))
			}
		}
		c.JSON(http.StatusOK, nil)
		return
	}

	jsonError(c, http.StatusInternalServerError, errors.New("case not found"))
}
