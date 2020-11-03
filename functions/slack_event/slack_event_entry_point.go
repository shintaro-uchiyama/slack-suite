package slack_event

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	previousSecretManager "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	"github.com/slack-go/slack"

	"github.com/dgrijalva/jwt-go"

	"cloud.google.com/go/datastore"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/slack-go/slack/slackevents"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type Card struct {
	ID int
}

type Response struct {
	AccessToken string `json:"access_token"`
}

type CreateCardRequest struct {
	AssigneeIds  []int  `json:"assignee_ids"`
	Body         string `json:"body"`
	CategoryName string `json:"category_name"`
	EpicId       int    `json:"epic_id"`
	GithubIssue  int    `json:"github_issue"`
	LabelIds     []int  `json:"label_ids"`
	Points       int    `json:"points"`
	Priority     int    `json:"priority"`
	ProjectId    int    `json:"project_id"`
	SprintId     int    `json:"sprint_id"`
	Title        string `json:"title"`
	WorkspaceId  int    `json:"workspace_id"`
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

	secretManagerClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("secretmanager client initialize error: %w", err))
		return err
	}

	zubePrivateKeySecret, err := secretManagerClient.AccessSecretVersion(ctx, &previousSecretManager.AccessSecretVersionRequest{
		Name: "projects/759555709793/secrets/zube-private-key/versions/latest",
	})
	if err != nil {
		log.Fatal(fmt.Errorf("get secret error: %w", err))
		return err
	}

	clientID := "83f007e2-1928-11eb-ac84-c7f4c49e7e6f"
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(10 * time.Hour).Unix(),
		Issuer:    clientID,
	})

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(zubePrivateKeySecret.Payload.Data)
	if err != nil {
		log.Fatal(fmt.Errorf("load signKey error: %w", err))
		return err
	}

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		log.Fatal(fmt.Errorf("get token error: %w", err))
		return err
	}

	httpClient := &http.Client{}
	httpReq, err := http.NewRequest("POST", "https://zube.io/api/users/tokens", nil)
	if err != nil {
		log.Fatal(fmt.Errorf("http req: %w", err))
		return err
	}
	httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	httpReq.Header.Add("X-Client-ID", clientID)
	httpReq.Header.Add("Accept", "application/json")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		log.Fatal(fmt.Errorf("client do: %w", err))
		return err
	}

	bodyB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("body byte: %w", err))
		return err
	}

	var response Response
	if err := json.Unmarshal(bodyB, &response); err != nil {
		log.Fatal(fmt.Errorf("unmarshal error: %w", err))
		return err
	}

	slackAccessTokenRequest := &previousSecretManager.AccessSecretVersionRequest{
		Name: "projects/759555709793/secrets/slack-access-token/versions/latest",
	}
	slackAccessToken, err := secretManagerClient.AccessSecretVersion(ctx, slackAccessTokenRequest)
	if err != nil {
		log.Fatal(fmt.Errorf("fetch slack signing secret error: %w", err))
		return err
	}

	api := slack.New(string(slackAccessToken.Payload.Data))
	conversationHistory, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{
		ChannelID: reactionAddedEvent.Item.Channel,
		Inclusive: true,
		Latest:    reactionAddedEvent.Item.Timestamp,
		Limit:     1,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("fetch conversation history error: %w", err))
		return err
	}

	httpClient = &http.Client{}
	body := CreateCardRequest{
		ProjectId: 25535,
		Title:     "test",
		Body:      conversationHistory.Messages[0].Text,
	}
	requestByte, _ := json.Marshal(body)

	httpReq, err = http.NewRequest("POST", "https://zube.io/api/cards", bytes.NewReader(requestByte))
	if err != nil {
		log.Fatal(fmt.Errorf("http req: %w", err))
		return err
	}
	httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", response.AccessToken))
	httpReq.Header.Add("X-Client-ID", "83f007e2-1928-11eb-ac84-c7f4c49e7e6f")
	httpReq.Header.Add("Content-Type", "application/json")

	resp, err = httpClient.Do(httpReq)
	if err != nil {
		log.Fatal(fmt.Errorf("client do: %w", err))
		return err
	}

	bodyB, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(fmt.Errorf("body byte: %w", err))
		os.Exit(0)
	}
	log.Println(fmt.Sprintf("prodject list: %+v", string(bodyB)))

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
