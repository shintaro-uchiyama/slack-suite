package application

import (
	"errors"
	"fmt"

	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"

	"github.com/slack-go/slack/slackevents"
)

type TaskApplication struct {
	projectService    ProjectServiceInterface
	taskService       TaskServiceInterface
	projectRepository domain.ProjectDataStoreInterface
	taskRepository    domain.TaskDataStoreInterface
	labelRepository   domain.LabelDataStoreInterface
	slack             domain.SlackInterface
	zube              domain.ZubeInterface
}

func NewTaskApplication(
	projectService ProjectServiceInterface,
	taskService TaskServiceInterface,
	projectRepository domain.ProjectDataStoreInterface,
	taskRepository domain.TaskDataStoreInterface,
	labelRepository domain.LabelDataStoreInterface,
	slack domain.SlackInterface,
	zube domain.ZubeInterface,
) *TaskApplication {
	return &TaskApplication{
		projectService:    projectService,
		taskService:       taskService,
		projectRepository: projectRepository,
		taskRepository:    taskRepository,
		labelRepository:   labelRepository,
		slack:             slack,
		zube:              zube,
	}
}

func (t TaskApplication) Create(reactionAddedEvent slackevents.ReactionAddedEvent) error {
	project, err := t.projectService.GetProjectByChannel(reactionAddedEvent.Item.Channel)
	if err != nil {
		return fmt.Errorf("can't create project by channel: %w", err)
	}

	task, err := project.CreateTask(t.labelRepository, reactionAddedEvent.Item.Timestamp, reactionAddedEvent.Reaction)
	if err != nil {
		return fmt.Errorf("can't create task: %w", err)
	}

	isExist, err := t.taskService.IsExist(task)
	if err != nil {
		return fmt.Errorf("is exist error %s", err)
	}
	if isExist {
		return fmt.Errorf("timeStamp %s already created", task.Timestamp())
	}

	_, err = t.taskService.Store(task)
	if err != nil {
		return fmt.Errorf("task persistence error %s", err)
	}
	return nil
}

func (t TaskApplication) Delete(reactionRemovedEvent slackevents.ReactionRemovedEvent) error {
	project, err := t.projectService.GetProjectByChannel(reactionRemovedEvent.Item.Channel)
	if err != nil {
		return fmt.Errorf("can't create project by channel: %w", err)
	}

	task, err := project.GetTaskByTimestamp(t.labelRepository, t.taskRepository, reactionRemovedEvent.Item.Timestamp, reactionRemovedEvent.Reaction)
	if err != nil {
		return fmt.Errorf("get task erro %s", err)
	}
	isExist, err := t.taskService.IsExist(task)
	if err != nil {
		return fmt.Errorf("is exist error %s", err)
	}
	if !isExist {
		return errors.New(fmt.Sprintf("timeStamp %s not exist", reactionRemovedEvent.Item.Timestamp))
	}

	err = project.DeleteTask(t.taskRepository, t.zube, task)
	if err != nil {
		return fmt.Errorf("delete task error %s", err)
	}
	return nil
}
