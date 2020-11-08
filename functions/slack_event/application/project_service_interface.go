package application

import (
	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"
)

var _ ProjectServiceInterface = (*domain.ProjectService)(nil)

type ProjectServiceInterface interface {
	GetByChannel(channel string) (*domain.Project, error)
}
