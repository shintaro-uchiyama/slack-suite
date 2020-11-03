package application

import (
	"fmt"
	"io"
	"net/http"

	"github.com/slack-go/slack/slackevents"
)

type VerifyApplication struct {
	secretManager SecretManagerInterface
	slackEvent    SlackEventInterface
}

func NewVerifyApplication(secretManager SecretManagerInterface, slackEvent SlackEventInterface) *VerifyApplication {
	return &VerifyApplication{
		secretManager: secretManager,
		slackEvent:    slackEvent,
	}
}

func (a VerifyApplication) Verify(header http.Header, body io.ReadCloser) ([]byte, error) {
	slackSigningSecret, err := a.secretManager.GetSecret("slack-signing-secret")
	if err != nil {
		return nil, fmt.Errorf("fetch slack signing secret error: %w", err)
	}

	bodyByte, err := a.slackEvent.Verify(header, body, slackSigningSecret)
	if err != nil {
		return nil, fmt.Errorf("slack secret verifier error: %w", err)
	}

	return bodyByte, nil
}

func (a VerifyApplication) ParseEvent(body []byte) (slackevents.EventsAPIEvent, error) {
	eventsAPIEvent, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
	return eventsAPIEvent, err
}
