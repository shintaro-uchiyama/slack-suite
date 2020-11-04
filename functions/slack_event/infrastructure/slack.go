package infrastructure

import (
	"errors"
	"fmt"
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type Slack struct {
	client *slack.Client
}

func NewSlack(slackAccessToken string) *Slack {
	return &Slack{
		client: slack.New(slackAccessToken),
	}
}

func (s Slack) GetMessage(item slackevents.Item) (string, error) {
	conversationHistory, err := s.client.GetConversationHistory(&slack.GetConversationHistoryParameters{
		ChannelID: item.Channel,
		Inclusive: true,
		Latest:    item.Timestamp,
		Limit:     1,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("fetch conversation history error: %w", err))
		return "", err
	}

	if len(conversationHistory.Messages) == 0 {
		return "", errors.New("message not found")
	}
	return conversationHistory.Messages[0].Text, nil
}
