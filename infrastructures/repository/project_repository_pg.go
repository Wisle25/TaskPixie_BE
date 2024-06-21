package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/task-pixie/applications/generator"
	"github.com/wisle25/task-pixie/domains/entity"
	"github.com/wisle25/task-pixie/domains/repository"
)

type ProjectRepositoryPG struct {
	idGenerator generator.IdGenerator
	db          *sql.DB
}

func NewProjectRepositoryPG(idGenerator generator.IdGenerator, db *sql.DB) repository.ProjectRepository {
	return &ProjectRepositoryPG{
		idGenerator,
		db,
	}
}

func (r *ProjectRepositoryPG) AddProject(payload *entity.ProjectPayload, ownerId string) string {
	id := r.idGenerator.Generate() // Assume you have a function to generate UUID

	query := `INSERT INTO projects(id, name, description, owner_id) VALUES ($1, $2, $3, $4) RETURNING id`
	var returnedId string
	err := r.db.QueryRow(query, id, payload.Name, payload.Description, ownerId).Scan(&returnedId)
	if err != nil {
		panic(fmt.Errorf("project_repo_pg_error: add project: %v", err))
	}
	return returnedId
}

func (r *ProjectRepositoryPG) GetProjectById(id string) *entity.Project {
	var project entity.Project
	query := `SELECT id, name, description, owner_id, created_at FROM projects WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&project.Id, &project.Name, &project.Description, &project.OwnerId, &project.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			panic(fiber.NewError(fiber.StatusNotFound, "Project not found!"))
		}
		panic(fmt.Errorf("project_repo_pg_error: get project by id: %v", err))
	}
	return &project
}

func (r *ProjectRepositoryPG) UpdateProjectById(id string, payload *entity.ProjectPayload) {
	query := `UPDATE projects SET name = $2, description = $3 WHERE id = $1`
	_, err := r.db.Exec(query, id, payload.Name, payload.Description)
	if err != nil {
		panic(fmt.Errorf("project_repo_pg_error: update project: %v", err))
	}
}

func (r *ProjectRepositoryPG) DeleteProjectById(id string) {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		panic(fmt.Errorf("project_repo_pg_error: delete project: %v", err))
	}
}

func (r *ProjectRepositoryPG) GetProjectsByOwner(ownerId string) []entity.Project {
	var projects []entity.Project
	query := `SELECT id, name, description, owner_id, created_at FROM projects WHERE owner_id = $1`
	rows, err := r.db.Query(query, ownerId)
	if err != nil {
		panic(fmt.Errorf("project_repo_pg_error: get projects by owner: %v", err))
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	for rows.Next() {
		var project entity.Project
		err := rows.Scan(&project.Id, &project.Name, &project.Description, &project.OwnerId, &project.CreatedAt)
		if err != nil {
			panic(fmt.Errorf("project_repo_pg_error: scan project: %v", err))
		}
		projects = append(projects, project)
	}
	return projects
}
