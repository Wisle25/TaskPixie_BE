package repository

import "github.com/wisle25/task-pixie/domains/entity"

// ProjectRepository defines methods for interacting with the project-related data in the database.
type ProjectRepository interface {
	AddProject(payload *entity.ProjectPayload, ownerId string) string
	GetProjectById(id string) *entity.Project
	UpdateProjectById(id string, payload *entity.ProjectPayload)
	DeleteProjectById(id string)
	GetProjectsByOwner(ownerId string) []entity.Project
}
