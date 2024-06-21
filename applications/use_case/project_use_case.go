package use_case

import (
	"github.com/wisle25/task-pixie/domains/entity"
	"github.com/wisle25/task-pixie/domains/repository"
)

// ProjectUseCase handles the business logic for project operations.
type ProjectUseCase struct {
	projectRepository repository.ProjectRepository
}

func NewProjectUseCase(projectRepository repository.ProjectRepository) *ProjectUseCase {
	return &ProjectUseCase{projectRepository: projectRepository}
}

func (uc *ProjectUseCase) ExecuteAddProject(payload *entity.ProjectPayload, ownerId string) string {
	return uc.projectRepository.AddProject(payload, ownerId)
}

func (uc *ProjectUseCase) ExecuteGetProjectById(id string) *entity.Project {
	return uc.projectRepository.GetProjectById(id)
}

func (uc *ProjectUseCase) ExecuteUpdateProjectById(id string, payload *entity.ProjectPayload) {
	uc.projectRepository.UpdateProjectById(id, payload)
}

func (uc *ProjectUseCase) ExecuteDeleteProjectById(id string) {
	uc.projectRepository.DeleteProjectById(id)
}

func (uc *ProjectUseCase) ExecuteGetProjectsByOwner(ownerId string) []entity.Project {
	return uc.projectRepository.GetProjectsByOwner(ownerId)
}
