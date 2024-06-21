package use_case

import (
	"github.com/wisle25/task-pixie/domains/entity"
	"github.com/wisle25/task-pixie/domains/repository"
)

type TaskUseCase struct {
	taskRepository repository.TaskRepository
}

func NewTaskUseCase(taskRepository repository.TaskRepository) *TaskUseCase {
	return &TaskUseCase{
		taskRepository: taskRepository,
	}
}

func (uc *TaskUseCase) ExecuteAddTask(task *entity.Task) string {
	return uc.taskRepository.AddTask(task)
}

func (uc *TaskUseCase) ExecuteGetTaskById(id string) *entity.Task {
	return uc.taskRepository.GetTaskById(id)
}

func (uc *TaskUseCase) ExecuteGetTasksByUser(userId string) []entity.Task {
	return uc.taskRepository.GetTasksByUser(userId)
}

func (uc *TaskUseCase) ExecuteGetTasksByProject(projectId string) []entity.Task {
	return uc.taskRepository.GetTasksByProject(projectId)
}

func (uc *TaskUseCase) ExecuteUpdateTask(task *entity.Task) string {
	return uc.taskRepository.UpdateTask(task)
}

func (uc *TaskUseCase) ExecuteDeleteTask(id string) {
	uc.taskRepository.DeleteTask(id)
}
