package application

import (
	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"
	"github.com/slack-go/slack/slackevents"
)

var _ TaskServiceInterface = (*domain.TaskService)(nil)

type TaskServiceInterface interface {
	IsExist(timeStamp string) (bool, error)
	Create(item slackevents.Item) error
	Delete(item slackevents.Item) error
}
