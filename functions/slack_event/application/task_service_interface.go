package application

import (
	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"
)

var _ TaskServiceInterface = (*domain.TaskService)(nil)

type TaskServiceInterface interface {
	IsExist(task domain.Task) (bool, error)
	Store(task domain.Task) (domain.Task, error)
}
