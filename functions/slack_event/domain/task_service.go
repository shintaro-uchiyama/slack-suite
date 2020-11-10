package domain

import (
	"fmt"
)

type TaskService struct {
	secretManager SecretManagerInterface
	slack         SlackInterface
	zube          ZubeInterface
	dataStore     DataStoreInterface
}

func NewTaskService(secretManager SecretManagerInterface, slack SlackInterface, zube ZubeInterface, dataStore DataStoreInterface) *TaskService {
	return &TaskService{
		secretManager: secretManager,
		slack:         slack,
		zube:          zube,
		dataStore:     dataStore,
	}
}

func (s TaskService) IsExist(task Task) (bool, error) {
	task, err := s.dataStore.Get(task.Timestamp())
	if err != nil {
		return false, fmt.Errorf("datastore get error: %w", err)
	}
	if task.CardID() == 0 {
		return false, nil
	}
	return true, nil
}
