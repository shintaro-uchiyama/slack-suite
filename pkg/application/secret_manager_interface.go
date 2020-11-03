package application

import "github.com/shintaro-uchiyama/pkg/infrastructure"

var _ secretManagerInterface = (*infrastructure.SecretManager)(nil)

type secretManagerInterface interface {
	GetSecret(secretName string) (string, error)
}
