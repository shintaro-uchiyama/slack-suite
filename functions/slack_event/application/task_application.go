package application

import (
	"errors"
	"fmt"

	"github.com/slack-go/slack/slackevents"
)

type TaskApplication struct {
	taskService TaskServiceInterface
}

func NewTaskApplication(taskService TaskServiceInterface) *TaskApplication {
	return &TaskApplication{
		taskService: taskService,
	}
}

func (t TaskApplication) Create(reactionAddedEvent slackevents.ReactionAddedEvent) error {
	isExist, err := t.taskService.IsExist(reactionAddedEvent.Item.Timestamp)
	if err != nil {
		return fmt.Errorf("is exist error %s", err)
	}
	if isExist {
		return errors.New(fmt.Sprintf("timeStamp %s already created", reactionAddedEvent.Item.Timestamp))
	}

	err = t.taskService.Create(reactionAddedEvent.Item)
	if err != nil {
		return fmt.Errorf("create task error %s", err)
	}
	return nil
}
