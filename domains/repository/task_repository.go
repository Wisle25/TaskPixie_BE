package repository

import "github.com/wisle25/task-pixie/domains/entity"

// TaskRepository defines methods for interacting with the task-related data in the database.
type TaskRepository interface {
	AddTask(task *entity.Task) string
	GetTaskById(id string) *entity.Task
	GetTasksByUser(userId string) []entity.Task
	GetTasksByProject(projectId string) []entity.Task
	UpdateTask(task *entity.Task) string
	DeleteTask(id string)
}
