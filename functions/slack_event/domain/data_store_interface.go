package domain

import (
	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/infrastructure"
)

var _ DataStoreInterface = (*infrastructure.DataStore)(nil)

type DataStoreInterface interface {
	Create(timeStamp string, cardID int) error
	Get(timeStamp string) (*infrastructure.Task, error)
}
