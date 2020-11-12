package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"

	"cloud.google.com/go/datastore"
)

var _ domain.LabelDataStoreInterface = (*LabelDataStore)(nil)

const labelDataStoreKey = "Label"

type LabelDataStore struct {
	client *datastore.Client
}

func NewLabelDataStore() (*LabelDataStore, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		return nil, fmt.Errorf("project datastore NewClient error: %w", err)
	}
	return &LabelDataStore{
		client: client,
	}, nil
}

type LabelEntity struct {
	LabelID int `datastore:",noindex"`
}

func (d LabelDataStore) GetByReaction(channel string, reaction string) (domain.Label, error) {
	ctx := context.Background()
	projectKey := datastore.NameKey(projectDataStoreKey, channel, nil)
	labelKey := datastore.NameKey(labelDataStoreKey, reaction, projectKey)
	var label LabelEntity
	err := d.client.Get(ctx, labelKey, &label)
	if errors.Is(err, datastore.ErrNoSuchEntity) {
		return domain.Label{}, nil
	} else if err != nil {
		return domain.Label{}, fmt.Errorf("get datastore error: %w", err)
	}
	return *domain.NewLabel(reaction, label.LabelID), nil
}
