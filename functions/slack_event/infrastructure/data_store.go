package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"

	"cloud.google.com/go/datastore"
)

var _ domain.DataStoreInterface = (*DataStore)(nil)

type DataStore struct {
	client *datastore.Client
	key    string
}

func NewDataStore() (*DataStore, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		return nil, fmt.Errorf("datastore NewClient error: %w", err)
	}
	return &DataStore{
		client: client,
		key:    "Task",
	}, nil
}

type Task struct {
	CardID int    `datastore:",noindex"`
	Title  string `datastore:",noindex"`
	Body   string `datastore:","`
}

func (d DataStore) Create(domainTask domain.Task) error {
	ctx := context.Background()
	projectKey := datastore.NameKey(d.key, domainTask.Project().Channel(), nil)
	key := datastore.NameKey(d.key, domainTask.Timestamp(), projectKey)
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

func (d DataStore) Delete(timestamp string) error {
	ctx := context.Background()
	key := datastore.NameKey(d.key, timestamp, nil)
	err := d.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("delete datastore error: %w", err)
	}
	return nil
}

func (d DataStore) Get(timestamp string) (domain.Task, error) {
	ctx := context.Background()
	key := datastore.NameKey(d.key, timestamp, nil)
	var task Task
	err := d.client.Get(ctx, key, &task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get datastore error: %w", err)
	}

	domainTask := domain.NewTask(domain.Project{}, timestamp, task.Title, task.Body, task.CardID)
	return *domainTask, nil
}
