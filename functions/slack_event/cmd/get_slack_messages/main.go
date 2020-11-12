package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/infrastructure"
)

func injectDependencies() (*infrastructure.Slack, error) {
	secretManager, err := infrastructure.NewSecretManager()
	if err != nil {
		return nil, fmt.Errorf("NewSecretManager error: %w", err)
	}
	slackAccessToken, err := secretManager.GetSecret("slack-access-token")
	if err != nil {
		return nil, fmt.Errorf("get zube private key secret error: %w", err)
	}
	slack := infrastructure.NewSlack(string(slackAccessToken))
	return slack, nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		logrus.Error(fmt.Sprintf("channel and timestamp args required bud %d", len(flag.Args())))
		os.Exit(1)
	}

	slack, err := injectDependencies()
	if err != nil {
		log.Fatal(fmt.Errorf("inject dependencies error: %w", err))
	}

	message, err := slack.GetMessage(flag.Arg(0), flag.Arg(1))
	if err != nil {
		log.Fatal(fmt.Errorf("get message error: %w", err))
	}
	log.Println(fmt.Sprintf("messages %+v, %+v", message.Body(), message.Title()))
}
