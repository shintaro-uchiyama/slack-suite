package application

import (
	"io"
	"net/http"

	"github.com/shintaro-uchiyama/pkg/infrastructure"
)

var _ SlackEventInterface = (*infrastructure.EventSlack)(nil)

type SlackEventInterface interface {
	Verify(header http.Header, body io.ReadCloser, slackSigningSecret string) ([]byte, error)
}
