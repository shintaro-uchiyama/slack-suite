package presentation

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"

	"cloud.google.com/go/pubsub"
	"github.com/slack-go/slack/slackevents"
)

type SlackEventHandler struct {
	taskApplication TaskApplicationInterface
}

func NewSlackEventHandler(taskApplication TaskApplicationInterface) *SlackEventHandler {
	return &SlackEventHandler{
		taskApplication: taskApplication,
	}
}

func (h SlackEventHandler) Create(ctx context.Context, m pubsub.Message) error {
	var reactionAddedEvent slackevents.ReactionAddedEvent
	if err := json.Unmarshal(m.Data, &reactionAddedEvent); err != nil {
		return fmt.Errorf("unmarshal pubsub message error: %w", err)
	}
	logrus.Debug(fmt.Sprintf("request add event: %+v", reactionAddedEvent))

	err := h.taskApplication.Create(reactionAddedEvent)
	if err != nil {
		return fmt.Errorf("create task error: %w", err)
	}
	return nil
}

func (h SlackEventHandler) Delete(ctx context.Context, m pubsub.Message) error {
	var reactionRemovedEvent slackevents.ReactionRemovedEvent
	if err := json.Unmarshal(m.Data, &reactionRemovedEvent); err != nil {
		return fmt.Errorf("unmarshal pubsub message error: %w", err)
	}
	logrus.Debug(fmt.Sprintf("request remove event: %+v", reactionRemovedEvent))

	err := h.taskApplication.Delete(reactionRemovedEvent)
	if err != nil {
		return fmt.Errorf("delete task error: %w", err)
	}
	return nil
}
