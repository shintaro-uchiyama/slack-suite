package domain

type Task struct {
	timeStamp string
	title     string
}

func NewTask(timeStamp string) *Task {
	return &Task{
		timeStamp: timeStamp,
	}
}
