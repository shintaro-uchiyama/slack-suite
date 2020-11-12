package domain

import "fmt"

type ProjectService struct {
	projectRepository ProjectDataStoreInterface
}

func NewProjectService(projectRepository ProjectDataStoreInterface) *ProjectService {
	return &ProjectService{
		projectRepository: projectRepository,
	}
}

func (s ProjectService) GetProjectByChannel(channel string) (Project, error) {
	project, err := s.projectRepository.GetByChannel(channel)
	if err != nil {
		return Project{}, fmt.Errorf("can't get project by channel: %w", err)
	}
	return project, nil
}
