package application

import (
	"fmt"
)

type TaskApplication struct {
	pubSub pubSubInterface
}

func NewTaskApplication(pubSub pubSubInterface) *TaskApplication {
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
