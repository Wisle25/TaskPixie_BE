package use_case

import (
	"github.com/wisle25/task-pixie/applications/validation"
	"github.com/wisle25/task-pixie/domains/entity"
	"github.com/wisle25/task-pixie/domains/repository"
)

// ProjectUseCase handles the business logic for project operations.
type ProjectUseCase struct {
	projectRepository repository.ProjectRepository
	validator         validation.ValidateProject
}

func NewProjectUseCase(
	projectRepository repository.ProjectRepository,
	validator validation.ValidateProject,
) *ProjectUseCase {
	return &ProjectUseCase{
		projectRepository: projectRepository,
		validator:         validator,
	}
}

// ExecuteAddProject handles the creation of a new project.
func (uc *ProjectUseCase) ExecuteAddProject(payload *entity.ProjectPayload, ownerId string) string {
	uc.validator.ValidatePayload(payload)
	return uc.projectRepository.AddProject(payload, ownerId)
}

// ExecuteGetProjectById retrieves a project by its ID.
func (uc *ProjectUseCase) ExecuteGetProjectById(id string) *entity.Project {
	return uc.projectRepository.GetProjectById(id)
}

func (uc *ProjectUseCase) ExecuteGetProjectMembers(id string) []entity.User {
	return uc.projectRepository.GetProjectMembers(id)
}

// ExecuteUpdateProjectById updates a project by its ID.
func (uc *ProjectUseCase) ExecuteUpdateProjectById(id string, payload *entity.ProjectPayload) {
	uc.validator.ValidatePayload(payload)
	uc.projectRepository.UpdateProjectById(id, payload)
}

// ExecuteDeleteProjectById deletes a project by its ID.
func (uc *ProjectUseCase) ExecuteDeleteProjectById(id string) {
	uc.projectRepository.DeleteProjectById(id)
}

// ExecuteGetProjects retrieves projects by owner or members.
func (uc *ProjectUseCase) ExecuteGetProjects(userId string) []entity.PreviewProject {
	ownerProjects := uc.projectRepository.GetProjectsByOwner(userId)
	memberProjects := uc.projectRepository.GetProjectsByMember(userId)

	projectMap := make(map[string]entity.PreviewProject)

	// Add owner projects to map
	for _, project := range ownerProjects {
		projectMap[project.Id] = entity.PreviewProject{
			Id:    project.Id,
			Title: project.Title,
		}
	}

	// Add member projects to map if not already added
	for _, project := range memberProjects {
		if _, exists := projectMap[project.Id]; !exists {
			projectMap[project.Id] = entity.PreviewProject{
				Id:    project.Id,
				Title: project.Title,
			}
		}
	}

	var previewProjects []entity.PreviewProject
	for _, project := range projectMap {
		previewProjects = append(previewProjects, project)
	}

	return previewProjects
}
