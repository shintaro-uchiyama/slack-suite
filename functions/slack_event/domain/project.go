package domain

import (
	"fmt"
)

type Project struct {
	id      int
	channel string
	tasks   []Task
}

func NewProject(id int, channel string) *Project {
	return &Project{
		id:      id,
		channel: channel,
	}
}

func (p Project) ID() int {
	return p.id
}

func (p Project) Channel() string {
	return p.channel
}

func (p Project) CreateTask(timeStamp string) (Task, error) {
	task := NewTask(p, timeStamp, "", "", 0)
	p.tasks = append(p.tasks, *task)
	return *task, nil
}

func (p *Project) GetTaskByTimestamp(taskRepository DataStoreInterface, timeStamp string) (Task, error) {
	task, err := taskRepository.Get(timeStamp)
	if err != nil {
		return Task{}, fmt.Errorf("get datastore error: %w", err)
	}
	return task, nil
}

func (p Project) Delete(taskRepository DataStoreInterface, zube ZubeInterface, cardID int, timestamp string) error {
	err := zube.Delete(cardID)
	if err != nil {
		return fmt.Errorf("delete zube card error: %w", err)
	}

	err = taskRepository.Delete(timestamp)
	if err != nil {
		return fmt.Errorf("delete datastore task error: %w", err)
	}

	return nil
}
