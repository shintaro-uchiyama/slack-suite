package application

import (
	"fmt"
)

type TaskApplication struct {
	pubSub PubSubInterface
}

func NewTaskApplication(pubSub PubSubInterface) *TaskApplication {
	return &TaskApplication{
		pubSub: pubSub,
	}
}

func (a TaskApplication) CallCreate(message []byte) error {
	err := a.pubSub.Publish("slack-event", message)
	if err != nil {
		return fmt.Errorf("topinc publish error: %w", err)
	}

	return nil
}

func (a TaskApplication) CallDelete(message []byte) error {
	err := a.pubSub.Publish("delete-task", message)
	if err != nil {
		return fmt.Errorf("delete-task topic publish error: %w", err)
	}

	return nil
}
