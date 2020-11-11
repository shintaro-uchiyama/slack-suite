package domain

type Label struct {
	reaction string
	id       int
}

func NewLabel(reaction string, id int) *Label {
	return &Label{
		reaction: reaction,
		id:       id,
	}
}

func (l Label) ID() int {
	return l.id
}
