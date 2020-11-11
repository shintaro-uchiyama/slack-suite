package infrastructure

import (
	"context"
	"fmt"
	"os"

	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"

	"cloud.google.com/go/datastore"
)

var _ domain.ProjectDataStoreInterface = (*ProjectDataStore)(nil)

const projectDataStoreKey = "Project"

type ProjectDataStore struct {
	client *datastore.Client
}

func NewProjectDataStore() (*ProjectDataStore, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		return nil, fmt.Errorf("project datastore NewClient error: %w", err)
	}
	return &ProjectDataStore{
		client: client,
	}, nil
}

type ProjectEntity struct {
	ProjectID int `datastore:",noindex"`
}

func (d ProjectDataStore) GetByChannel(channel string) (domain.Project, error) {
	ctx := context.Background()
	key := datastore.NameKey(projectDataStoreKey, channel, nil)
	var project ProjectEntity
	err := d.client.Get(ctx, key, &project)
	if err != nil {
		return domain.Project{}, fmt.Errorf("get datastore error: %w", err)
	}
	return *domain.NewProject(project.ProjectID, channel), nil
}
