package presentation

import (
	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/application"
	"github.com/slack-go/slack/slackevents"
)

var _ TaskApplicationInterface = (*application.TaskApplication)(nil)

type TaskApplicationInterface interface {
	Create(reactionAddedEvent slackevents.ReactionAddedEvent) error
}
