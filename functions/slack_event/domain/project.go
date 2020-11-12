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

func (p *Project) CreateTask(labelRepository LabelDataStoreInterface, timeStamp string, reaction string) (Task, error) {
	label, err := labelRepository.GetByReaction(p.channel, reaction)
	if err != nil {
		return Task{}, fmt.Errorf("can't get label by reaction: %w", err)
	}
	task := NewTask(*p, timeStamp, "", "", 0, label.ID())
	p.tasks = append(p.tasks, *task)
	return *task, nil
}

func (p *Project) GetTaskByTimestamp(labelRepository LabelDataStoreInterface, taskRepository TaskDataStoreInterface, timeStamp string, reaction string) (Task, error) {
	_, err := labelRepository.GetByReaction(p.channel, reaction)
	if err != nil {
		return Task{}, fmt.Errorf("can't get label by reaction: %w", err)
	}

	task, err := taskRepository.Get(p.channel, timeStamp)
	if err != nil {
		return Task{}, fmt.Errorf("get datastore error: %w", err)
	}
	task.SetProject(*p)
	return *task, nil
}

func (p Project) DeleteTask(taskRepository TaskDataStoreInterface, zube ZubeInterface, task Task) error {
	err := zube.Delete(task.CardID())
	if err != nil {
		return fmt.Errorf("delete zube card error: %w", err)
	}

	err = taskRepository.Delete(p.channel, task.Timestamp())
	if err != nil {
		return fmt.Errorf("delete datastore task error: %w", err)
	}

	return nil
}
