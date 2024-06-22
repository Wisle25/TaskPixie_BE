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
	db          *sql.DB
	idGenerator generator.IdGenerator
}

func NewProjectRepositoryPG(db *sql.DB, idGenerator generator.IdGenerator) repository.ProjectRepository {
	return &ProjectRepositoryPG{
		db:          db,
		idGenerator: idGenerator,
	}
}

func (r *ProjectRepositoryPG) AddProject(payload *entity.ProjectPayload, ownerId string) string {
	// Create ID
	id := r.idGenerator.Generate()

	// Insert project
	query := `INSERT INTO 
    			projects(id, title, detail, priority, status, owner_id) 
			  VALUES
			      ($1, $2, $3, $4, $5, $6)
			  RETURNING id`

	var returnedId string
	err := r.db.QueryRow(
		query,
		id,
		payload.Title,
		payload.Detail,
		payload.Priority,
		payload.Status,
		ownerId,
	).Scan(&returnedId)

	if err != nil {
		panic(fmt.Errorf("project_repo_pg_error: add project: %v", err))
	}

	// Insert members
	for _, memberId := range payload.MembersId {
		query := `INSERT INTO project_members(project_id, user_id) VALUES ($1, $2)`
		_, err := r.db.Exec(query, returnedId, memberId)
		if err != nil {
			panic(fmt.Errorf("project_repo_pg_error: add project member: %v", err))
		}
	}

	return returnedId
}

func (r *ProjectRepositoryPG) GetProjectById(id string) *entity.Project {
	var project entity.Project
	var memberUsername string

	// Query project details
	query := `SELECT id, title, detail, priority, status, created_at, updated_at FROM projects WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&project.Id,
		&project.Title,
		&project.Detail,
		&project.Priority,
		&project.Status,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			panic(fiber.NewError(fiber.StatusNotFound, "Project not found!"))
		}
		panic(fmt.Errorf("project_repo_pg_error: get project by id: %v", err))
	}

	// Query project members
	query = `SELECT u.username FROM project_members pm JOIN users u ON pm.user_id = u.id WHERE pm.project_id = $1`
	rows, err := r.db.Query(query, id)
	if err != nil {
		panic(fmt.Errorf("project_repo_pg_error: get project members: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&memberUsername)
		if err != nil {
			panic(fmt.Errorf("project_repo_pg_error: scan project member: %v", err))
		}
		project.MembersUsername = append(project.MembersUsername, memberUsername)
	}

	return &project
}

func (r *ProjectRepositoryPG) GetProjectMembers(id string) []entity.User {
	query := `
        SELECT
            u.id,
            u.username
        FROM project_members pm
        INNER JOIN users u ON pm.user_id = u.id
        WHERE pm.project_id = $1`

	rows, err := r.db.Query(query, id)
	if err != nil {
		panic(fmt.Errorf("project_repo_pg_error: get project members: %v", err))
	}
	defer rows.Close()

	var members []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.Id, &user.Username); err != nil {
			panic(fmt.Errorf("project_repo_pg_error: scan project member: %v", err))
		}
		members = append(members, user)
	}

	return members
}

func (r *ProjectRepositoryPG) UpdateProjectById(id string, payload *entity.ProjectPayload) {
	// Update project details
	query := `
		UPDATE projects 
		SET title = $1, detail = $2, priority = $3, status = $4, updated_at = NOW()
		WHERE id = $5`

	_, err := r.db.Exec(query, payload.Title, payload.Detail, payload.Priority, payload.Status, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			panic(fiber.NewError(fiber.StatusNotFound, "Project not found!"))
		}
		panic(fmt.Errorf("project_repo_pg_error: update project: %v", err))
	}

	// Remove existing members
	query = `DELETE FROM project_members WHERE project_id = $1`
	_, err = r.db.Exec(query, id)
	if err != nil {
		panic(fmt.Errorf("project_repo_pg_error: delete project members: %v", err))
	}

	// Insert new members
	for _, memberId := range payload.MembersId {
		query := `INSERT INTO project_members(project_id, user_id) VALUES ($1, $2)`
		_, err := r.db.Exec(query, id, memberId)
		if err != nil {
			panic(fmt.Errorf("project_repo_pg_error: add project member: %v", err))
		}
	}
}

func (r *ProjectRepositoryPG) DeleteProjectById(id string) {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.Exec(query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			panic(fiber.NewError(fiber.StatusNotFound, "Project not found!"))
		}
		panic(fmt.Errorf("project_repo_pg_error: delete project: %v", err))
	}

	// Delete project members
	query = `DELETE FROM project_members WHERE project_id = $1`
	_, err = r.db.Exec(query, id)
	if err != nil {
		panic(fmt.Errorf("project_repo_pg_error: delete project members: %v", err))
	}
}

func (r *ProjectRepositoryPG) GetProjectsByOwner(ownerId string) []entity.PreviewProject {
	var projects []entity.PreviewProject

	// Query projects by owner
	query := `SELECT id, title FROM projects WHERE owner_id = $1`
	rows, err := r.db.Query(query, ownerId)
	if err != nil {
		panic(fmt.Errorf("project_repo_pg_error: get projects by owner: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		var project entity.PreviewProject
		err := rows.Scan(&project.Id, &project.Title)
		if err != nil {
			panic(fmt.Errorf("project_repo_pg_error: scan project: %v", err))
		}
		projects = append(projects, project)
	}

	return projects
}

func (r *ProjectRepositoryPG) GetProjectsByMember(memberId string) []entity.PreviewProject {
	var projects []entity.PreviewProject

	// Query projects by member
	query := `
		SELECT p.id, p.title
		FROM projects p
		JOIN project_members pm ON p.id = pm.project_id
		WHERE pm.user_id = $1`
	rows, err := r.db.Query(query, memberId)
	if err != nil {
		panic(fmt.Errorf("project_repo_pg_error: get projects by member: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		var project entity.PreviewProject
		err := rows.Scan(&project.Id, &project.Title)
		if err != nil {
			panic(fmt.Errorf("project_repo_pg_error: scan project: %v", err))
		}
		projects = append(projects, project)
	}

	return projects
}
