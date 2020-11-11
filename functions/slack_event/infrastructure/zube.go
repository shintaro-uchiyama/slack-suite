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

	"github.com/sirupsen/logrus"

	"github.com/shintaro-uchiyama/slack-suite/functions/slack_event/domain"

	"github.com/dgrijalva/jwt-go"
)

var _ domain.ZubeInterface = (*Zube)(nil)

type Zube struct {
	accessToken string
	clientID    string
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

	return &Zube{
		accessToken: response.AccessToken,
		clientID:    os.Getenv("CLIENT_ID"),
	}, nil
}

type CreateCardRequest struct {
	ProjectId    int    `json:"project_id"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	LabelIds     []int  `json:"label_ids"`
	CategoryName string `json:"category_name"`
	EpicId       int    `json:"epic_id"`
	GithubIssue  int    `json:"github_issue"`
	AssigneeIds  []int  `json:"assignee_ids"`
	Points       int    `json:"points"`
	Priority     int    `json:"priority"`
	SprintId     int    `json:"sprint_id"`
	WorkspaceId  int    `json:"workspace_id"`
}

type CreateCardResponse struct {
	ID int `json:"id"`
}

func (z Zube) Create(task domain.Task) (int, error) {
	logrus.Info(fmt.Sprintf("labels %+v: ", task.Labels()))
	requestByte, err := json.Marshal(CreateCardRequest{
		ProjectId:   task.Project().ID(),
		WorkspaceId: task.Project().WorkspaceID(),
		Title:       task.Title(),
		Body:        task.Body(),
		LabelIds:    task.Labels(),
	})
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

	httpClient := &http.Client{}
	resp, _ := httpClient.Do(httpReq)
	bodyByte, _ := ioutil.ReadAll(resp.Body)

	var response CreateCardResponse
	if err := json.Unmarshal(bodyByte, &response); err != nil {
		return 0, fmt.Errorf("unmarshal zube response error: %w", err)
	}
	return response.ID, nil
}

func (z Zube) Delete(cardID int) error {
	httpClient := &http.Client{}
	httpReq, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("https://zube.io/api/cards/%s/archive", strconv.Itoa(cardID)),
		nil,
	)

	if err != nil {
		return fmt.Errorf("zube archive http request error: %w", err)
	}
	httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", z.accessToken))
	httpReq.Header.Add("X-Client-ID", z.clientID)
	httpReq.Header.Add("Content-Type", "application/json")

	_, err = httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("zube archive http request error: %w", err)
	}
	return nil
}

type Workspace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Project struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Workspaces []Workspace `json:"workspaces"`
}

type ResponseBody struct {
	Data []Project `json:"data"`
}

func (z Zube) GetProjects() ([]Project, error) {
	httpClient := &http.Client{}
	httpReq, err := http.NewRequest("GET", "https://zube.io/api/projects", nil)

	if err != nil {
		return nil, fmt.Errorf("zube get projects http request error: %w", err)
	}
	httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", z.accessToken))
	httpReq.Header.Add("X-Client-ID", z.clientID)
	httpReq.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("convert zube project id from string to int error: %w", err)
	}
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read zube get projects http body error: %w", err)
	}

	var response ResponseBody
	if err := json.Unmarshal(bodyByte, &response); err != nil {
		return nil, fmt.Errorf("unmarshal zube response error: %w", err)
	}
	return response.Data, nil
}

type Label struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LabelsResponse struct {
	Data []Label `json:"data"`
}

func (z Zube) GetLabels(zubeProjectID string) (interface{}, error) {
	httpClient := &http.Client{}
	httpReq, err := http.NewRequest("GET", fmt.Sprintf("https://zube.io/api/projects/%s/labels", zubeProjectID), nil)

	if err != nil {
		return nil, fmt.Errorf("zube get labels http request error: %w", err)
	}
	httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", z.accessToken))
	httpReq.Header.Add("X-Client-ID", z.clientID)
	httpReq.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request error: %w", err)
	}
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read zube labels http body error: %w", err)
	}

	var response LabelsResponse
	if err := json.Unmarshal(bodyByte, &response); err != nil {
		return nil, fmt.Errorf("unmarshal zube response error: %w", err)
	}
	return response, nil
}
