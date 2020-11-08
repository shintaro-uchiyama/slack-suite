package application

import (
	"errors"
	"fmt"

	"github.com/slack-go/slack/slackevents"
)

type TaskApplication struct {
	taskService    TaskServiceInterface
	projectService ProjectServiceInterface
}

func NewTaskApplication(taskService TaskServiceInterface, projectService ProjectServiceInterface) *TaskApplication {
	return &TaskApplication{
		taskService:    taskService,
		projectService: projectService,
	}
}

func (t TaskApplication) Create(reactionAddedEvent slackevents.ReactionAddedEvent) error {
	project, err := t.projectService.GetByChannel(reactionAddedEvent.Item.Channel)
	if err != nil {
		return fmt.Errorf("get project error %s", err)
	}
	if project == nil {
		return fmt.Errorf("project not found %s", err)
	}

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

func (t TaskApplication) Delete(reactionRemovedEvent slackevents.ReactionRemovedEvent) error {
	isExist, err := t.taskService.IsExist(reactionRemovedEvent.Item.Timestamp)
	if err != nil {
		return fmt.Errorf("is exist error %s", err)
	}
	if !isExist {
		return errors.New(fmt.Sprintf("timeStamp %s not exist", reactionRemovedEvent.Item.Timestamp))
	}

	err = t.taskService.Delete(reactionRemovedEvent.Item)
	if err != nil {
		return fmt.Errorf("delete task error %s", err)
	}
	return nil
}
