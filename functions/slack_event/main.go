package slack_event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/datastore"
	"github.com/slack-go/slack/slackevents"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type Card struct {
	ID int
}

func SlackEventEntryPoint(ctx context.Context, m PubSubMessage) error {
	var reactionAddedEvent slackevents.ReactionAddedEvent
	if err := json.Unmarshal(m.Data, &reactionAddedEvent); err != nil {
		log.Printf(fmt.Errorf("unmarshal error: %w", err).Error())
		return err
	}
	log.Printf("reactionAddedEvent, %+v!", reactionAddedEvent)

	var err error
	client, err := datastore.NewClient(ctx, "uchiyama-sandbox")
	if err != nil {
		log.Fatal(err)
		return err
	}

	taskKey := datastore.NameKey("Card", reactionAddedEvent.Item.Timestamp, nil)
	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		// We first check that there is no entity stored with the given key.
		var empty Card
		if err := tx.Get(taskKey, &empty); err != datastore.ErrNoSuchEntity {
			return err
		}
		// If there was no matching entity, store it now.
		_, err := tx.Put(taskKey, &Card{
			ID: 1,
		})
		return err
	})

	if err != nil {
		log.Fatal(fmt.Errorf("upsert error: %w", err))
		return err
	}

	return nil
}
