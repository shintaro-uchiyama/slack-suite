package application

import (
	"io"
	"net/http"

	"github.com/shintaro-uchiyama/pkg/infrastructure"
)

var _ slackEventInterface = (*infrastructure.EventSlack)(nil)

type slackEventInterface interface {
	Verify(header http.Header, body io.ReadCloser, slackSigningSecret string) ([]byte, error)
}
