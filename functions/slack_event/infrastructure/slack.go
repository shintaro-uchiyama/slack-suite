package infrastructure

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

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
	conversationReplies, _, _, err := s.client.GetConversationReplies(
		&slack.GetConversationRepliesParameters{
			ChannelID: channel,
			Inclusive: true,
			Timestamp: timestamp,
		},
	)
	if err != nil {
		return domain.SlackMessage{}, fmt.Errorf("fetch conversation history error: %w", err)
	}

	if len(conversationReplies) == 0 {
		return domain.SlackMessage{}, errors.New("message not found")
	}

	conversationReply := conversationReplies[0]
	logrus.Info(conversationReply)
	text := conversationReply.Text
	title, body := text, text
	index := strings.Index(text, "\n")
	if index > -1 {
		title = text[:index]
	}

	var linkUrl string
	if conversationReply.ThreadTimestamp != "" {
		linkUrl = fmt.Sprintf(
			"%s/%s/p%s?thread_ts=%s&cid=%s",
			os.Getenv("SLACK_URL"),
			channel,
			strings.Replace(timestamp, ".", "", -1),
			conversationReply.ThreadTimestamp,
			channel,
		)
	} else {
		linkUrl = fmt.Sprintf(
			"%s/%s/p%s",
			os.Getenv("SLACK_URL"),
			channel,
			strings.Replace(timestamp, ".", "", -1),
		)
	}
	body = fmt.Sprintf("%s \n %s", body, linkUrl)
	return *domain.NewSlackMessage(title, body), nil
}
