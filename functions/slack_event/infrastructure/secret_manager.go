package infrastructure

import (
	"context"
	"fmt"
	"os"

	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"

	secretManager "cloud.google.com/go/secretmanager/apiv1"
	previousSecretManager "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var _ domain.SecretManagerInterface = (*SecretManager)(nil)

type SecretManager struct {
	client        *secretManager.Client
	projectNumber string
}

func NewSecretManager() (*SecretManager, error) {
	ctx := context.Background()
	client, err := secretManager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("secret manager new client error: %w", err)
	}
	projectNumber := os.Getenv("PROJECT_NUMBER")
	return &SecretManager{
		client:        client,
		projectNumber: projectNumber,
	}, nil
}

func (s SecretManager) GetSecret(secretName string) ([]byte, error) {
	ctx := context.Background()
	secret, err := s.client.AccessSecretVersion(ctx, &previousSecretManager.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", s.projectNumber, secretName),
	})
	if err != nil {
		return nil, fmt.Errorf("access secret error: %w", err)
	}
	return secret.Payload.Data, nil
}
