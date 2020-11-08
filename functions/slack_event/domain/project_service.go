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

func (s ProjectService) GetByChannel(channel string) (*Project, error) {
	projectEntity, err := s.projectDataStore.GetByChannel(channel)
	if err != nil {
		return nil, fmt.Errorf("get project ID error %s", err)
	}
	if projectEntity == nil {
		return nil, fmt.Errorf("project not found by channel %s", err)
	}
	project := NewProject(projectEntity.ProjectID, channel)
	return project, nil
}
