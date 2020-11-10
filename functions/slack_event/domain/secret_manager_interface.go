package domain

type SecretManagerInterface interface {
	GetSecret(secretName string) ([]byte, error)
}
