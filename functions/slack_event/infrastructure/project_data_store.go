package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"

	"cloud.google.com/go/datastore"
)

var _ domain.ProjectDataStoreInterface = (*ProjectDataStore)(nil)

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
	ProjectID int `datastore:",noindex"`
}

func (d ProjectDataStore) GetByChannel(channel string) (domain.Project, error) {
	ctx := context.Background()
	key := datastore.NameKey(d.kind, channel, nil)
	var project ProjectEntity
	err := d.client.Get(ctx, key, &project)
	if errors.Is(err, datastore.ErrNoSuchEntity) {
		return domain.Project{}, datastore.ErrNoSuchEntity
	} else if err != nil {
		return domain.Project{}, fmt.Errorf("get datastore error: %w", err)
	}
	return *domain.NewProject(project.ProjectID, channel), nil
}
