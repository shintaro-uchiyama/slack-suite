package domain

type Task struct {
	channel   string
	timeStamp string
	title     string
}

func NewTask(channel string, timeStamp string) *Task {
	return &Task{
		channel:   channel,
		timeStamp: timeStamp,
	}
}
