package domain

type TaskService struct {
	taskRepository TaskDataStoreInterface
}

func NewTaskService(taskRepository TaskDataStoreInterface) *TaskService {
	return &TaskService{
		taskRepository: taskRepository,
	}
}

func (s TaskService) IsExist(task Task) (bool, error) {
	foundTask, err := s.taskRepository.Get(task.Project().Channel(), task.Timestamp())
	if err != nil {
		return false, nil
	}
	if foundTask == nil {
		return false, nil
	}
	return true, nil
}
