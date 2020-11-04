package slack_event

import (
	"context"
	"fmt"
	"os"

	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/presentation"

	"cloud.google.com/go/pubsub"
	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/application"
	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"
	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/infrastructure"
	"github.com/sirupsen/logrus"
)

func initLog() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
}

func injectDependencies() (*presentation.SlackEventHandler, error) {
	secretManager, err := infrastructure.NewSecretManager()
	if err != nil {
		return nil, fmt.Errorf("NewSecretManager: %w", err)
	}
	slackAccessToken, err := secretManager.GetSecret("slack-access-token")
	if err != nil {
		return nil, fmt.Errorf("get slack access token secret error: %w", err)
	}
	zubePrivateKey, err := secretManager.GetSecret("zube-private-key")
	if err != nil {
		return nil, fmt.Errorf("get zube private key secret error: %w", err)
	}

	slack := infrastructure.NewSlack(string(slackAccessToken))
	zube, err := infrastructure.NewZube(zubePrivateKey)
	if err != nil {
		return nil, fmt.Errorf("NewZube error: %w", err)
	}

	dataStore, err := infrastructure.NewDataStore()
	if err != nil {
		return nil, fmt.Errorf("NewDataStore error: %w", err)
	}

	taskApplication := presentation.NewSlackEventHandler(
		application.NewTaskApplication(
			domain.NewTaskService(secretManager, slack, zube, dataStore),
		))
	return taskApplication, nil
}

func SlackEventEntryPoint(ctx context.Context, m pubsub.Message) error {
	initLog()

	taskApplication, err := injectDependencies()
	if err != nil {
		return fmt.Errorf("inject dependencies error: %w", err)
	}

	err = taskApplication.Create(ctx, m)
	if err != nil {
		err = fmt.Errorf("create task error: %w", err)
		logrus.Errorf(err.Error())
		return err
	}
	return nil
}

func DeleteTaskEntryPoint(ctx context.Context, m pubsub.Message) error {
	initLog()

	taskApplication, err := injectDependencies()
	if err != nil {
		return fmt.Errorf("inject dependencies error: %w", err)
	}

	err = taskApplication.Delete(ctx, m)
	if err != nil {
		err = fmt.Errorf("create task error: %w", err)
		logrus.Errorf(err.Error())
		return err
	}
	return nil
}
