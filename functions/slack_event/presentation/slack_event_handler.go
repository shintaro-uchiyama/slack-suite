package presentation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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

var targetReactions = map[string]int{"zube": 0}

func (h SlackEventHandler) Create(ctx context.Context, m pubsub.Message) error {
	var reactionAddedEvent slackevents.ReactionAddedEvent
	if err := json.Unmarshal(m.Data, &reactionAddedEvent); err != nil {
		return fmt.Errorf("unmarshal pubsub message error: %w", err)
	}

	if _, ok := targetReactions[reactionAddedEvent.Reaction]; !ok {
		return errors.New(fmt.Sprintf("%s is not target reaction", reactionAddedEvent.Reaction))
	}
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

	if _, ok := targetReactions[reactionRemovedEvent.Reaction]; !ok {
		return errors.New(fmt.Sprintf("%s is not target reactoin", reactionRemovedEvent.Reaction))
	}
	err := h.taskApplication.Delete(reactionRemovedEvent)
	if err != nil {
		return fmt.Errorf("delete task error: %w", err)
	}
	return nil
}
