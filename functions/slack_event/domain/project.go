package domain

type Project struct {
	projectID int    `json:"project_id"`
	channel   string `json:"channel"`
}

func NewProject(projectID int, channel string) *Project {
	return &Project{
		projectID: projectID,
		channel:   channel,
	}
}

func (p Project) GetProjectID() int {
	return p.projectID
}
