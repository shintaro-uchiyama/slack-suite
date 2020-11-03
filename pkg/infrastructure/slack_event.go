package infrastructure

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/slack-go/slack"
)

type EventSlack struct{}

func NewEventSlack() *EventSlack {
	return &EventSlack{}
}

func (e EventSlack) Verify(header http.Header, body io.ReadCloser, slackSigningSecret string) ([]byte, error) {
	verifier, err := slack.NewSecretsVerifier(header, slackSigningSecret)
	if err != nil {
		return nil, fmt.Errorf("slack secret verifier error: %w", err)
	}

	bodyReader := io.TeeReader(body, &verifier)
	bodyByte, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, fmt.Errorf("read request body error: %w", err)
	}

	if err := verifier.Ensure(); err != nil {
		return nil, fmt.Errorf("ensure slack secret verifier error: %w", err)
	}
	return bodyByte, nil
}
