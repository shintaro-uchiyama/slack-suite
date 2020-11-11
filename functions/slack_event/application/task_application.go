package application

import (
	"errors"
	"fmt"

	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"

	"github.com/slack-go/slack/slackevents"
)

type TaskApplication struct {
	taskService       TaskServiceInterface
	projectRepository domain.ProjectDataStoreInterface
	taskRepository    domain.TaskDataStoreInterface
	labelRepository   domain.LabelDataStoreInterface
	slack             domain.SlackInterface
	zube              domain.ZubeInterface
}

func NewTaskApplication(
	taskService TaskServiceInterface,
	projectRepository domain.ProjectDataStoreInterface,
	taskRepository domain.TaskDataStoreInterface,
	labelRepository domain.LabelDataStoreInterface,
	slack domain.SlackInterface,
	zube domain.ZubeInterface,
) *TaskApplication {
	return &TaskApplication{
		taskService:       taskService,
		projectRepository: projectRepository,
		taskRepository:    taskRepository,
		labelRepository:   labelRepository,
		slack:             slack,
		zube:              zube,
	}
}

func (t TaskApplication) Create(reactionAddedEvent slackevents.ReactionAddedEvent) error {
	project, err := t.projectRepository.GetByChannel(reactionAddedEvent.Item.Channel)
	if err != nil {
		return fmt.Errorf("get project entity error %w", err)
	}

	task, err := project.CreateTask(reactionAddedEvent.Item.Timestamp)
	if err != nil {
		return fmt.Errorf("crete task error %w", err)
	}

	isExist, err := t.taskService.IsExist(task)
	if err != nil {
		return fmt.Errorf("is exist error %s", err)
	}
	if isExist {
		return errors.New(fmt.Sprintf("timeStamp %s already created", reactionAddedEvent.Item.Timestamp))
	}

	slackMessage, err := t.slack.GetMessage(reactionAddedEvent.Item.Channel, reactionAddedEvent.Item.Timestamp)
	if err != nil {
		return fmt.Errorf("get slack message error: %w", err)
	}
	task.SetTitle(slackMessage.Title())
	task.SetBody(slackMessage.Body())

	label, err := t.labelRepository.GetByReaction(reactionAddedEvent.Item.Channel, reactionAddedEvent.Reaction)
	if err != nil {
		return fmt.Errorf("get label datastore error: %w", err)
	}
	task.AddLabel(label.ID())

	cardID, err := t.zube.Create(task)
	if err != nil {
		return fmt.Errorf("create zube card error: %w", err)
	}
	task.SetCardID(cardID)

	err = t.taskRepository.Create(task)
	if err != nil {
		deleteErr := t.zube.Delete(cardID)
		if deleteErr != nil {
			return fmt.Errorf("delete zube card error: %w", err)
		}
		return fmt.Errorf("create datastore task error: %w", err)
	}

	return nil
}

func (t TaskApplication) Delete(reactionRemovedEvent slackevents.ReactionRemovedEvent) error {
	project, err := t.projectRepository.GetByChannel(reactionRemovedEvent.Item.Channel)
	if err != nil {
		return fmt.Errorf("get project entity error %w", err)
	}

	task, err := project.GetTaskByTimestamp(t.taskRepository, reactionRemovedEvent.Item.Timestamp)
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

	err = project.Delete(t.taskRepository, t.zube, task.CardID(), task.Timestamp())
	if err != nil {
		return fmt.Errorf("delete task error %s", err)
	}
	return nil
}
