package application

import "github.com/shintaro-uchiyama/pkg/infrastructure"

var _ SecretManagerInterface = (*infrastructure.SecretManager)(nil)

type SecretManagerInterface interface {
	GetSecret(secretName string) ([]byte, error)
}
