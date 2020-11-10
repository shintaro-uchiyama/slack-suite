package domain

type SlackInterface interface {
	GetMessage(channel string, timestamp string) (SlackMessage, error)
}
