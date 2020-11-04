package domain

import "github.com/shintaro-uchiyama/slack-suite/functions/slack_event/infrastructure"

var _ SecretManagerInterface = (*infrastructure.SecretManager)(nil)

type SecretManagerInterface interface {
	GetSecret(secretName string) ([]byte, error)
}
