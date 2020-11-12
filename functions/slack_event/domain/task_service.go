package domain

import (
	"fmt"
)

type TaskService struct {
	taskRepository TaskDataStoreInterface
	slack          SlackInterface
	zube           ZubeInterface
}

func NewTaskService(taskRepository TaskDataStoreInterface, slack SlackInterface, zube ZubeInterface) *TaskService {
	return &TaskService{
		taskRepository: taskRepository,
		slack:          slack,
		zube:           zube,
	}
}

func (s TaskService) IsExist(task Task) (bool, error) {
	foundTask, err := s.taskRepository.Get(task.Project().Channel(), task.Timestamp())
	if err != nil {
		return false, nil
	}
	if foundTask == nil {
		return false, nil
	}
	return true, nil
}

func (s TaskService) Store(task Task) (Task, error) {
	slackMessage, err := s.slack.GetMessage(task.Project().Channel(), task.Timestamp())
	if err != nil {
		return Task{}, fmt.Errorf("get slack message error: %w", err)
	}
	task.SetTitle(slackMessage.Title())
	task.SetBody(slackMessage.Body())

	cardID, err := s.zube.Create(task)
	if err != nil {
		return Task{}, fmt.Errorf("create zube card error: %w", err)
	}
	task.SetCardID(cardID)

	err = s.taskRepository.Create(task)
	if err != nil {
		deleteErr := s.zube.Delete(cardID)
		if deleteErr != nil {
			return Task{}, fmt.Errorf("delete zube card error: %w", deleteErr)
		}
		return Task{}, fmt.Errorf("create datastore task error: %w", err)
	}
	return task, nil
}
