package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/datastore"
)

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
	CardID int
}

func (d DataStore) Create(timeStamp string, cardID int) error {
	ctx := context.Background()
	key := datastore.NameKey(d.key, timeStamp, nil)
	_, err := d.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var task Task
		err := tx.Get(key, &task)
		if err != nil && !errors.Is(err, datastore.ErrNoSuchEntity) {
			return fmt.Errorf("get task from datastore error: %w", err)
		}

		_, err = tx.Put(key, &Task{
			CardID: cardID,
		})
		if err != nil {
			return fmt.Errorf("pubt task error: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("datastore trnsaction error: %w", err)
	}
	return nil
}

func (d DataStore) Delete(timeStamp string) error {
	ctx := context.Background()
	key := datastore.NameKey(d.key, timeStamp, nil)
	err := d.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("delete datastore error: %w", err)
	}
	return nil
}

func (d DataStore) Get(timeStamp string) (*Task, error) {
	ctx := context.Background()
	key := datastore.NameKey(d.key, timeStamp, nil)
	var task Task
	err := d.client.Get(ctx, key, &task)
	if errors.Is(err, datastore.ErrNoSuchEntity) {
		return &task, nil
	} else if err != nil {
		return nil, fmt.Errorf("get datastore error: %w", err)
	}
	return &task, nil
}
