package domain

import (
	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/infrastructure"
)

var _ ProjectDataStoreInterface = (*infrastructure.ProjectDataStore)(nil)

type ProjectDataStoreInterface interface {
	GetByChannel(channel string) (*infrastructure.ProjectEntity, error)
}
