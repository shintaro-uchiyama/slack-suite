package domain

import (
	"fmt"

	"github.com/slack-go/slack/slackevents"
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

func (s TaskService) IsExist(timeStamp string) (bool, error) {
	task, err := s.dataStore.Get(timeStamp)
	if err != nil {
		return false, fmt.Errorf("datastore get error: %w", err)
	}
	if task.CardID == 0 {
		return false, nil
	}
	return true, nil
}

func (s TaskService) Create(item slackevents.Item) error {
	text, err := s.slack.GetMessageText(item)
	if err != nil {
		return fmt.Errorf("get slack message error: %w", err)
	}

	item.Message.Text = text
	cardID, err := s.zube.Create(item)
	if err != nil {
		return fmt.Errorf("create zube card error: %w", err)
	}

	err = s.dataStore.Create(item.Timestamp, cardID)
	if err != nil {
		deleteErr := s.zube.Delete(cardID)
		if deleteErr != nil {
			return fmt.Errorf("delete zube card error: %w", err)
		}
		return fmt.Errorf("create datastore task error: %w", err)
	}

	return nil
}

func (s TaskService) Delete(item slackevents.Item) error {
	task, err := s.dataStore.Get(item.Timestamp)
	if err != nil {
		return fmt.Errorf("get datastore error: %w", err)
	}

	err = s.zube.Delete(task.CardID)
	if err != nil {
		return fmt.Errorf("delete zube card error: %w", err)
	}

	err = s.dataStore.Delete(item.Timestamp)
	if err != nil {
		return fmt.Errorf("delete datastore task error: %w", err)
	}

	return nil
}
