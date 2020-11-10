package infrastructure

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"

	"github.com/slack-go/slack"
)

var _ domain.SlackInterface = (*Slack)(nil)

type Slack struct {
	client *slack.Client
}

func NewSlack(slackAccessToken string) *Slack {
	return &Slack{
		client: slack.New(slackAccessToken),
	}
}

type SlackMessage struct {
	Title string
	Body  string
}

func (s Slack) GetMessage(channel string, timestamp string) (domain.SlackMessage, error) {
	conversationHistory, err := s.client.GetConversationHistory(
		&slack.GetConversationHistoryParameters{
			ChannelID: channel,
			Inclusive: true,
			Latest:    timestamp,
			Limit:     1,
		},
	)
	if err != nil {
		return domain.SlackMessage{}, fmt.Errorf("fetch conversation history error: %w", err)
	}

	if len(conversationHistory.Messages) == 0 {
		return domain.SlackMessage{}, errors.New("message not found")
	}

	text := conversationHistory.Messages[0].Text
	title, body := text, text
	index := strings.Index(text, "\n")
	if index > -1 {
		title = text[:index]
	}

	linkUrl := fmt.Sprintf(
		"%s/%s/p%s",
		os.Getenv("SLACK_URL"),
		channel,
		strings.Replace(channel, ".", "", -1),
	)
	body = fmt.Sprintf("%s \n %s", body, linkUrl)
	return *domain.NewSlackMessage(title, body), nil
}
