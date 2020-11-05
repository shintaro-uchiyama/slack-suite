package domain

import (
	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/infrastructure"
	"github.com/slack-go/slack/slackevents"
)

var _ ZubeInterface = (*infrastructure.Zube)(nil)

type ZubeInterface interface {
	Create(item slackevents.Item) (int, error)
	Delete(cardID int) error
}
