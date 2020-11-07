package application

import (
	"encoding/json"
	"fmt"

	"github.com/slack-go/slack/slackevents"
)

type TaskApplication struct {
	pubSub PubSubInterface
}

func NewTaskApplication(pubSub PubSubInterface) *TaskApplication {
	return &TaskApplication{
		pubSub: pubSub,
	}
}

func (a TaskApplication) CallCreate(event *slackevents.ReactionAddedEvent) error {
	messageByte, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	err = a.pubSub.Publish("create-task", messageByte)
	if err != nil {
		return fmt.Errorf("topinc publish error: %w", err)
	}

	return nil
}

func (a TaskApplication) CallDelete(event *slackevents.ReactionRemovedEvent) error {
	messageByte, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	err = a.pubSub.Publish("delete-task", messageByte)
	if err != nil {
		return fmt.Errorf("delete-task topic publish error: %w", err)
	}

	return nil
}
