package domain

type SlackMessage struct {
	title string
	body  string
}

func NewSlackMessage(title string, body string) *SlackMessage {
	return &SlackMessage{
		title: title,
		body:  body,
	}
}

func (s SlackMessage) Title() string {
	return s.title
}

func (s SlackMessage) Body() string {
	return s.body
}
