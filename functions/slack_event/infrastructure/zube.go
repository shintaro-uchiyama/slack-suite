package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Zube struct {
	accessToken   string
	zubeProjectID int
	clientID      string
}

type ZubeTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func NewZube(zubePrivateKey []byte) (*Zube, error) {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(zubePrivateKey)
	if err != nil {
		return nil, fmt.Errorf("load signKey error: %w", err)
	}

	clientID := os.Getenv("CLIENT_ID")
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		Issuer:    clientID,
	})
	refreshToken, err := token.SignedString(signKey)
	if err != nil {
		return nil, fmt.Errorf("get token string error: %w", err)
	}

	httpClient := &http.Client{}
	httpReq, err := http.NewRequest("POST", "https://zube.io/api/users/tokens", nil)
	if err != nil {
		return nil, fmt.Errorf("zube token http request error: %w", err)
	}
	httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", refreshToken))
	httpReq.Header.Add("X-Client-ID", clientID)
	httpReq.Header.Add("Accept", "application/json")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("access token erquest error: %w", err)
	}

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read access token response error: %w", err)
	}

	var response ZubeTokenResponse
	if err := json.Unmarshal(bodyByte, &response); err != nil {
		return nil, fmt.Errorf("unmarshal access token response error: %w", err)
	}

	zubeProjectIDString := os.Getenv("ZUBE_PROJECT_ID")
	zubeProjectID, err := strconv.Atoi(zubeProjectIDString)
	if err != nil {
		return nil, fmt.Errorf("convert zube project id from string to int error: %w", err)
	}

	return &Zube{
		accessToken:   response.AccessToken,
		zubeProjectID: zubeProjectID,
		clientID:      os.Getenv("CLIENT_ID"),
	}, nil
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

type CreateCardResponse struct {
	Number int `json:"number"`
}

func (z Zube) Create(title string, body string) (int, error) {
	httpClient := &http.Client{}
	createCardRequest := CreateCardRequest{
		ProjectId: z.zubeProjectID,
		Title:     title,
		Body:      body,
		LabelIds:  []int{272338},
	}
	requestByte, err := json.Marshal(createCardRequest)
	if err != nil {
		return 0, fmt.Errorf("createCardRequest error: %w", err)
	}

	httpReq, err := http.NewRequest("POST", "https://zube.io/api/cards", bytes.NewReader(requestByte))

	if err != nil {
		return 0, fmt.Errorf("zube create http request error: %w", err)
	}
	httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", z.accessToken))
	httpReq.Header.Add("X-Client-ID", z.clientID)
	httpReq.Header.Add("Content-Type", "application/json")

	resp, _ := httpClient.Do(httpReq)
	bodyByte, _ := ioutil.ReadAll(resp.Body)

	var response CreateCardResponse
	if err := json.Unmarshal(bodyByte, &response); err != nil {
		return 0, fmt.Errorf("unmarshal zube response error: %w", err)
	}
	return response.Number, nil
}
