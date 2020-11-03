//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/shintaro-uchiyama/pkg/application"
	"github.com/shintaro-uchiyama/pkg/infrastructure"
	"github.com/shintaro-uchiyama/pkg/presentation"
)

func InitializeEvent() (*presentation.EventHandler, error) {
	wire.Build(
		presentation.NewEventHandler,
		application.NewVerifyApplication,
		wire.Bind(new(presentation.VerifyApplicationInterface), new(*application.VerifyApplication)),
		infrastructure.NewSecretManager,
		wire.Bind(new(application.SecretManagerInterface), new(*infrastructure.SecretManager)),
		infrastructure.NewEventSlack,
		wire.Bind(new(application.SlackEventInterface), new(*infrastructure.EventSlack)),
		application.NewTaskApplication,
		wire.Bind(new(presentation.TaskApplicationInterface), new(*application.TaskApplication)),
		infrastructure.NewPubSub,
		wire.Bind(new(application.PubSubInterface), new(*infrastructure.PubSub)),
	)
	return nil, nil
}
