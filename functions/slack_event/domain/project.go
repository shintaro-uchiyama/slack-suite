package domain

import (
	"fmt"
)

type Project struct {
	id          int
	workspaceID int
	channel     string
	tasks       []Task
}

func NewProject(id int, channel string, workspaceID int) *Project {
	return &Project{
		id:          id,
		channel:     channel,
		workspaceID: workspaceID,
	}
}

func (p Project) ID() int {
	return p.id
}

func (p Project) WorkspaceID() int {
	return p.workspaceID
}

func (p Project) Channel() string {
	return p.channel
}

func (p Project) CreateTask(timeStamp string) (Task, error) {
	task := NewTask(p, timeStamp, "", "", 0)
	p.tasks = append(p.tasks, *task)
	return *task, nil
}

func (p *Project) GetTaskByTimestamp(taskRepository TaskDataStoreInterface, timeStamp string) (Task, error) {
	task, err := taskRepository.Get(p.channel, timeStamp)
	if err != nil {
		return Task{}, fmt.Errorf("get datastore error: %w", err)
	}
	task.SetProject(*p)
	return *task, nil
}

func (p Project) Delete(taskRepository TaskDataStoreInterface, zube ZubeInterface, cardID int, timestamp string) error {
	err := zube.Delete(cardID)
	if err != nil {
		return fmt.Errorf("delete zube card error: %w", err)
	}

	err = taskRepository.Delete(p.channel, timestamp)
	if err != nil {
		return fmt.Errorf("delete datastore task error: %w", err)
	}

	return nil
}
