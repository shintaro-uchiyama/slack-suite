package application

import "github.com/shintaro-uchiyama/pkg/infrastructure"

var _ pubSubInterface = (*infrastructure.PubSub)(nil)

type pubSubInterface interface {
	Publish(topicName string, message []byte) error
}
