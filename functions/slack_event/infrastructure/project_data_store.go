package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/datastore"
)

type ProjectDataStore struct {
	client *datastore.Client
	kind   string
}

func NewProjectDataStore() (*ProjectDataStore, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		return nil, fmt.Errorf("project datastore NewClient error: %w", err)
	}
	return &ProjectDataStore{
		client: client,
		kind:   "Project",
	}, nil
}

type ProjectEntity struct {
	Project int `datastore:",noindex"`
}

func (d ProjectDataStore) GetByChannel(channel string) (*ProjectEntity, error) {
	ctx := context.Background()
	key := datastore.NameKey(d.kind, channel, nil)
	var project ProjectEntity
	err := d.client.Get(ctx, key, &project)
	if errors.Is(err, datastore.ErrNoSuchEntity) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("get datastore error: %w", err)
	}
	return &project, nil
}
