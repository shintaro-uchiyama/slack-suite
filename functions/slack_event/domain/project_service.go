package domain

import (
	"fmt"
)

type ProjectService struct {
	projectDataStore ProjectDataStoreInterface
}

func NewProjectService(projectDataStore ProjectDataStoreInterface) *ProjectService {
	return &ProjectService{
		projectDataStore: projectDataStore,
	}
}

func (s ProjectService) GetByChannel(channel string) (Project, error) {
	project, err := s.projectDataStore.GetByChannel(channel)
	if err != nil {
		return Project{}, fmt.Errorf("get project ID error %s", err)
	}
	return project, nil
}
