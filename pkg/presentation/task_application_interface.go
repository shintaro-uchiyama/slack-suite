package presentation

import "github.com/shintaro-uchiyama/pkg/application"

var _ TaskApplicationInterface = (*application.TaskApplication)(nil)

type TaskApplicationInterface interface {
	CallCreate(message []byte) error
	CallDelete(message []byte) error
}
