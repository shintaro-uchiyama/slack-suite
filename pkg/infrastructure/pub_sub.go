package infrastructure

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
)

type PubSub struct {
	client        *pubsub.Client
	projectNumber string
}

func NewPubSub() (*PubSub, error) {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub client error: %w", err)
	}
	projectNumber := os.Getenv("PROJECT_NUMBER")
	return &PubSub{
		client:        client,
		projectNumber: projectNumber,
	}, nil
}

func (p PubSub) Publish(topicName string, message []byte) error {
	topic := p.client.Topic(topicName)
	ctx := context.Background()
	if _, err := topic.Publish(ctx, &pubsub.Message{
		Data: message,
	}).Get(ctx); err != nil {
		return fmt.Errorf("could not publish message: %w", err)
	}
	return nil
}
