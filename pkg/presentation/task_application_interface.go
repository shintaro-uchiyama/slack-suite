package presentation

import "github.com/shintaro-uchiyama/pkg/application"

var _ taskApplicationInterface = (*application.TaskApplication)(nil)

type taskApplicationInterface interface {
	CallCreate(message []byte) error
}
