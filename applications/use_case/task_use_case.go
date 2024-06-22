package use_case

import (
	"github.com/wisle25/task-pixie/applications/validation"
	"github.com/wisle25/task-pixie/domains/entity"
	"github.com/wisle25/task-pixie/domains/repository"
	"log"
)

// TaskUseCase handles the business logic for task operations.
type TaskUseCase struct {
	taskRepository repository.TaskRepository
	validator      validation.ValidateTask
}

func NewTaskUseCase(
	taskRepository repository.TaskRepository,
	validator validation.ValidateTask,
) *TaskUseCase {
	return &TaskUseCase{
		taskRepository: taskRepository,
		validator:      validator,
	}
}

// ExecuteAddTask handles the creation of a new task.
func (uc *TaskUseCase) ExecuteAddTask(payload *entity.TaskPayload, ownerId string) string {
	uc.validator.ValidatePayload(payload)
	return uc.taskRepository.AddTask(payload, ownerId)
}

// ExecuteGetTaskById retrieves a task by its ID.
func (uc *TaskUseCase) ExecuteGetTaskById(id string) *entity.Task {
	return uc.taskRepository.GetTaskById(id)
}

func (uc *TaskUseCase) ExecuteGetTasksByProjects(projectId string) []entity.PreviewTask {
	return uc.taskRepository.GetTasksByProjects(projectId)
}

// ExecuteUpdateTaskById updates a task by its ID.
func (uc *TaskUseCase) ExecuteUpdateTaskById(id string, payload *entity.TaskPayload) {
	uc.validator.ValidatePayload(payload)
	uc.taskRepository.UpdateTaskById(id, payload)
}

// ExecuteDeleteTaskById deletes a task by its ID.
func (uc *TaskUseCase) ExecuteDeleteTaskById(id string) {
	uc.taskRepository.DeleteTaskById(id)
}

// ExecuteGetTasks retrieves tasks by owner or assignees.
func (uc *TaskUseCase) ExecuteGetTasks(userId string) []entity.PreviewTask {
	log.Printf("User Use Case: %s", userId)
	ownerTasks := uc.taskRepository.GetTasksByOwner(userId)
	assigneeTasks := uc.taskRepository.GetTasksByAssignedUser(userId)

	taskMap := make(map[string]entity.PreviewTask)

	// Add owner tasks to map
	for _, task := range ownerTasks {
		taskMap[task.ID] = entity.PreviewTask{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Priority:    task.Priority,
			Status:      task.Status,
			Project:     task.Project,
		}
	}

	// Add assignee tasks to map if not already added
	for _, task := range assigneeTasks {
		if _, exists := taskMap[task.ID]; !exists {
			taskMap[task.ID] = entity.PreviewTask{
				ID:          task.ID,
				Title:       task.Title,
				Description: task.Description,
				Priority:    task.Priority,
				Status:      task.Status,
				Project:     task.Project,
			}
		}
	}

	var previewTasks []entity.PreviewTask
	for _, task := range taskMap {
		previewTasks = append(previewTasks, task)
	}

	return previewTasks
}
