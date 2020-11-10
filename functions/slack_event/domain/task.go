package domain

type Task struct {
	project   Project
	timestamp string
	title     string
	body      string
	cardID    int
}

func NewTask(project Project, timestamp string, title string, body string, cardID int) *Task {
	return &Task{
		project:   project,
		timestamp: timestamp,
		title:     title,
		body:      body,
		cardID:    cardID,
	}
}

func (t Task) Project() Project {
	return t.project
}

func (t Task) Timestamp() string {
	return t.timestamp
}

func (t *Task) SetTitle(title string) {
	t.title = title
}

func (t Task) SetBody(body string) {
	t.body = body
}

func (t Task) SetCardID(cardID int) {
	t.cardID = cardID
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Body() string {
	return t.body
}

func (t Task) CardID() int {
	return t.cardID
}
