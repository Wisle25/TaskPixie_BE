package repository

import "github.com/wisle25/task-pixie/domains/entity"

// TaskRepository defines methods for interacting with the task-related data in the database.
type TaskRepository interface {
	AddTask(payload *entity.TaskPayload, ownerId string) string
	GetTaskById(id string) *entity.Task
	UpdateTaskById(id string, payload *entity.TaskPayload)
	DeleteTaskById(id string)
	GetTasksByProjects(projectId string) []entity.PreviewTask
	GetTasksByOwner(ownerId string) []entity.PreviewTask
	GetTasksByAssignedUser(userId string) []entity.PreviewTask
}
