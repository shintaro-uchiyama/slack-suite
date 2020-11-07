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

func (a VerifyApplication) Verify(header http.Header, body io.ReadCloser) (*slackevents.EventsAPIEvent, []byte, error) {
	slackSigningSecret, err := a.secretManager.GetSecret("slack-signing-secret")
	if err != nil {
		return nil, nil, fmt.Errorf("fetch slack signing secret error: %w", err)
	}

	bodyByte, err := a.slackEvent.Verify(header, body, string(slackSigningSecret))
	if err != nil {
		return nil, nil, fmt.Errorf("slack secret verifier error: %w", err)
	}

	eventsAPIEvent, err := slackevents.ParseEvent(bodyByte, slackevents.OptionNoVerifyToken())
	if err != nil {
		return nil, nil, fmt.Errorf("slack event parse error: %w", err)
	}
	return &eventsAPIEvent, bodyByte, nil
}
