package presentation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/slack-go/slack/slackevents"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	verifyApplication VerifyApplicationInterface
	taskApplication   TaskApplicationInterface
}

func NewEventHandler(verifyApplication VerifyApplicationInterface, taskApplication TaskApplicationInterface) *EventHandler {
	return &EventHandler{
		verifyApplication: verifyApplication,
		taskApplication:   taskApplication,
	}
}

var targetReactions = map[string]int{"zube": 0}

func (h EventHandler) Create(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		_ = c.Error(fmt.Errorf("read request body error: %w", err)).SetType(gin.ErrorTypePrivate)
		return
	}

	apiEvent, err := h.verifyApplication.ParseEvent(body)
	if err != nil {
		_ = c.Error(fmt.Errorf("error found in verify: %w", err)).SetType(gin.ErrorTypePublic)
		return
	}

	switch apiEvent.Type {
	case slackevents.URLVerification:
		var challengeResponse *slackevents.ChallengeResponse
		err := json.Unmarshal(body, &challengeResponse)
		if err != nil {
			_ = c.Error(fmt.Errorf("slack url verification error: %w", err)).SetType(gin.ErrorTypePrivate)
			return
		}
		c.JSON(http.StatusOK, challengeResponse.Challenge)
	case slackevents.CallbackEvent:
		bodyByte, err := h.verifyApplication.Verify(c.Request.Header, c.Request.Body)
		if err != nil {
			_ = c.Error(fmt.Errorf("error found in verify: %w", err)).SetType(gin.ErrorTypePublic)
			return
		}

		slackEvent, err := h.verifyApplication.ParseEvent(bodyByte)
		if err != nil {
			_ = c.Error(fmt.Errorf("parse slack eventsAPI error: %w", err)).SetType(gin.ErrorTypePrivate)
			return
		}

		switch event := slackEvent.InnerEvent.Data.(type) {
		case *slackevents.ReactionAddedEvent:
			if _, ok := targetReactions[event.Reaction]; !ok {
				c.JSON(http.StatusOK, nil)
				return
			}

			messageByte, err := json.Marshal(event)
			if err != nil {
				_ = c.Error(fmt.Errorf("json marshal error: %w", err)).SetType(gin.ErrorTypePrivate)
				return
			}

			err = h.taskApplication.CallCreate(messageByte)
			if err != nil {
				_ = c.Error(fmt.Errorf("call create error: %w", err)).SetType(gin.ErrorTypePrivate)
				return
			}
			c.JSON(http.StatusOK, nil)
		case *slackevents.ReactionRemovedEvent:
			if _, ok := targetReactions[event.Reaction]; !ok {
				logrus.Info("not target remove reaction")
				c.JSON(http.StatusOK, nil)
				return
			}
			if event.Item.Channel != os.Getenv("CHANNEL_ID") {
				logrus.Info("not target channel")
				c.JSON(http.StatusOK, nil)
				return
			}

			messageByte, err := json.Marshal(event)
			if err != nil {
				_ = c.Error(fmt.Errorf("json marshal error: %w", err)).SetType(gin.ErrorTypePrivate)
				return
			}

			err = h.taskApplication.CallDelete(messageByte)
			if err != nil {
				_ = c.Error(fmt.Errorf("call delete error: %w", err)).SetType(gin.ErrorTypePrivate)
				return
			}
			c.JSON(http.StatusOK, nil)
		}
	default:
		_ = c.Error(
			errors.New(fmt.Sprintf("expected slack event not found, got %s", apiEvent.Type)),
		).SetType(gin.ErrorTypePublic)
		return
	}
}
