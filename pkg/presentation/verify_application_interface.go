package presentation

import (
	"io"
	"net/http"

	"github.com/shintaro-uchiyama/pkg/application"

	"github.com/slack-go/slack/slackevents"
)

var _ VerifyApplicationInterface = (*application.VerifyApplication)(nil)

type VerifyApplicationInterface interface {
	Verify(header http.Header, body io.ReadCloser) ([]byte, error)
	ParseEvent(body []byte) (slackevents.EventsAPIEvent, error)
}
