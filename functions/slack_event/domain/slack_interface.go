package domain

import (
	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/infrastructure"
	"github.com/slack-go/slack/slackevents"
)

var _ SlackInterface = (*infrastructure.Slack)(nil)

type SlackInterface interface {
	GetMessage(item slackevents.Item) (string, error)
}
