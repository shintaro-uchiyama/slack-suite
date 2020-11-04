package domain

import "github.com/shintaro-uchiyama/slack-suite/functions/slack_event/infrastructure"

var _ ZubeInterface = (*infrastructure.Zube)(nil)

type ZubeInterface interface {
	Create(title string, body string) (int, error)
	Delete(cardID int) error
}
