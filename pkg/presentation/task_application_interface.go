package presentation

import (
	"github.com/shintaro-uchiyama/pkg/application"
	"github.com/slack-go/slack/slackevents"
)

var _ TaskApplicationInterface = (*application.TaskApplication)(nil)

type TaskApplicationInterface interface {
	CallCreate(event *slackevents.ReactionAddedEvent) error
	CallDelete(event *slackevents.ReactionRemovedEvent) error
}
