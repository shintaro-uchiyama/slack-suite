package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"

	"cloud.google.com/go/datastore"
)

var _ domain.TaskDataStoreInterface = (*TaskDataStore)(nil)

const taskDataStoreKey = "Task"

type TaskDataStore struct {
	client *datastore.Client
}

func NewTaskDataStore() (*TaskDataStore, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		return nil, fmt.Errorf("datastore NewClient error: %w", err)
	}
	return &TaskDataStore{
		client: client,
	}, nil
}

type Task struct {
	CardID int    `datastore:",noindex"`
	Title  string `datastore:",noindex"`
	Body   string `datastore:""`
}

func (d TaskDataStore) Create(domainTask domain.Task) error {
	ctx := context.Background()
	projectKey := datastore.NameKey(projectDataStoreKey, domainTask.Project().Channel(), nil)
	key := datastore.NameKey(taskDataStoreKey, domainTask.Timestamp(), projectKey)
	_, err := d.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var task Task
		if err := tx.Get(key, &task); !errors.Is(err, datastore.ErrNoSuchEntity) {
			return fmt.Errorf("get task from datastore error: %w", err)
		}

		_, err := tx.Put(key, &Task{
			CardID: domainTask.CardID(),
			Title:  domainTask.Title(),
			Body:   domainTask.Body(),
		})
		if err != nil {
			return fmt.Errorf("put task error: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("datastore trnsaction error: %w", err)
	}
	return nil
}

func (d TaskDataStore) Delete(channel string, timestamp string) error {
	ctx := context.Background()
	projectKey := datastore.NameKey(projectDataStoreKey, channel, nil)
	key := datastore.NameKey(taskDataStoreKey, timestamp, projectKey)
	err := d.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("delete datastore error: %w", err)
	}
	return nil
}

func (d TaskDataStore) Get(channel string, timestamp string) (*domain.Task, error) {
	ctx := context.Background()
	projectKey := datastore.NameKey(projectDataStoreKey, channel, nil)
	key := datastore.NameKey(taskDataStoreKey, timestamp, projectKey)
	var task Task
	err := d.client.Get(ctx, key, &task)
	if errors.Is(err, datastore.ErrNoSuchEntity) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("get datastore error: %w", err)
	}

	domainTask := domain.NewTask(domain.Project{}, timestamp, task.Title, task.Body, task.CardID, 0)
	return domainTask, nil
}
