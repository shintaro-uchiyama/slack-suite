package application

import "github.com/shintaro-uchiyama/pkg/infrastructure"

var _ PubSubInterface = (*infrastructure.PubSub)(nil)

type PubSubInterface interface {
	Publish(topicName string, message []byte) error
}
