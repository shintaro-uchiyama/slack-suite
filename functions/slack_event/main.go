package slack_event

import (
	"context"
	"encoding/json"
	"log"

	"github.com/slack-go/slack/slackevents"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func SlackEventEntryPoint(ctx context.Context, m PubSubMessage) error {
	var obj slackevents.ReactionAddedEvent
	if err := json.Unmarshal(m.Data, &obj); err != nil {
		log.Printf("error decoding sakura response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("sakura response: %q", m.Data)
		return err
	}
	log.Printf("obj, %+v!", obj)

	return nil
}
