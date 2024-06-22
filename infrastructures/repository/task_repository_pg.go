package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/task-pixie/applications/generator"
	"github.com/wisle25/task-pixie/domains/entity"
	"github.com/wisle25/task-pixie/domains/repository"
	"github.com/wisle25/task-pixie/infrastructures/services"
	"log"
)

type TaskRepositoryPG struct {
	idGenerator generator.IdGenerator
	db          *sql.DB
}

func NewTaskRepositoryPG(idGenerator generator.IdGenerator, db *sql.DB) repository.TaskRepository {
	return &TaskRepositoryPG{
		idGenerator: idGenerator,
		db:          db,
	}
}

func (r *TaskRepositoryPG) AddTask(payload *entity.TaskPayload, ownerId string) string {
	// Create ID
	id := r.idGenerator.Generate()

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		panic(fmt.Errorf("task_repo_pg_error: begin transaction: %v", err))
	}

	// Defer a rollback in case anything fails
	defer tx.Rollback()

	// Base query for inserting task
	query := `INSERT INTO tasks (id, title, description, detail, priority, status, due_date, owner_id`
	args := []interface{}{id, payload.Title, payload.Description, payload.Detail, payload.Priority, payload.Status, payload.DueDate, ownerId}

	// Handle optional project_id
	if payload.ProjectId != "" {
		query += `, project_id`
		args = append(args, payload.ProjectId)
	}

	query += `) VALUES ($1, $2, $3, $4, $5, $6, $7, $8`
	if payload.ProjectId != "" {
		query += `, $9`
	}
	query += `) RETURNING id`

	var returnedId string
	err = tx.QueryRow(query, args...).Scan(&returnedId)
	if err != nil {
		panic(fmt.Errorf("task_repo_pg_error: add task: %v", err))
	}

	// Insert task assignments
	if len(payload.AssignedToId) > 0 {
		for _, userId := range payload.AssignedToId {
			assignmentQuery := `INSERT INTO task_assignments (task_id, user_id) VALUES ($1, $2)`
			_, err := tx.Exec(assignmentQuery, returnedId, userId)
			if err != nil {
				panic(fmt.Errorf("task_repo_pg_error: add task assignments: %v", err))
			}
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		panic(fmt.Errorf("task_repo_pg_error: commit transaction: %v", err))
	}

	return returnedId
}

func (r *TaskRepositoryPG) GetTaskById(id string) *entity.Task {
	var task entity.Task
	var projectId sql.NullString
	var project sql.NullString
	var assignedToUsernames []string

	// Query to get task details
	taskQuery := `SELECT t.id, t.title, t.description, t.detail, t.priority, t.status, p.id AS projectId, p.title as project, t.due_date, t.created_at, t.updated_at 
				  FROM tasks t
				  LEFT JOIN projects p ON t.project_id = p.id
				  WHERE t.id = $1`
	err := r.db.QueryRow(taskQuery, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Detail,
		&task.Priority,
		&task.Status,
		&projectId,
		&project,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	task.ProjectId = projectId.String
	task.Project = project.String

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			panic(fiber.NewError(fiber.StatusNotFound, "Task not found!"))
		}
		panic(fmt.Errorf("task_repo_pg_error: get task by id: %v", err))
	}

	// Query to get assigned usernames
	assignmentsQuery := `SELECT u.username 
						 FROM users u
						 JOIN task_assignments ta ON u.id = ta.user_id
						 WHERE ta.task_id = $1`
	rows, err := r.db.Query(assignmentsQuery, id)
	if err != nil {
		panic(fmt.Errorf("task_repo_pg_error: get task assignments: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			panic(fmt.Errorf("task_repo_pg_error: scan username: %v", err))
		}
		assignedToUsernames = append(assignedToUsernames, username)
	}

	task.AssignedToUsernames = assignedToUsernames

	return &task
}

func (r *TaskRepositoryPG) GetTasksByProjects(projectId string) []entity.PreviewTask {
	// Query
	query := `
			SELECT 
				t.id, t.title, t.description, t.priority, t.status, p.title as project
			FROM tasks t
			INNER JOIN projects p ON p.id = t.project_id
			WHERE t.project_id = $1`
	rows, err := r.db.Query(query, projectId)

	if err != nil {
		panic(fmt.Errorf("task_repo_pg_error: get tasks by project: %v", err))
	}

	return services.GetTableDB[entity.PreviewTask](rows)
}

func (r *TaskRepositoryPG) UpdateTaskById(id string, payload *entity.TaskPayload) {
	// Query to update task
	query := `UPDATE tasks SET title = $1, description = $2, detail = $3, priority = $4, status = $5, project_id = $6, due_date = $7, updated_at = NOW() 
			  WHERE id = $8`

	_, err := r.db.Exec(
		query,
		payload.Title,
		payload.Description,
		payload.Detail,
		payload.Priority,
		payload.Status,
		payload.ProjectId,
		payload.DueDate,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			panic(fiber.NewError(fiber.StatusNotFound, "Task not found!"))
		}
		panic(fmt.Errorf("task_repo_pg_error: update task: %v", err))
	}

	// Delete existing task assignments
	deleteQuery := `DELETE FROM task_assignments WHERE task_id = $1`
	_, err = r.db.Exec(deleteQuery, id)
	if err != nil {
		panic(fmt.Errorf("task_repo_pg_error: delete task assignments: %v", err))
	}

	// Insert new task assignments
	if len(payload.AssignedToId) > 0 {
		for _, userId := range payload.AssignedToId {
			assignmentQuery := `INSERT INTO task_assignments (task_id, user_id) VALUES ($1, $2)`
			_, err := r.db.Exec(assignmentQuery, id, userId)
			if err != nil {
				panic(fmt.Errorf("task_repo_pg_error: add task assignments: %v", err))
			}
		}
	}
}

func (r *TaskRepositoryPG) DeleteTaskById(id string) {
	// Query to delete task
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			panic(fiber.NewError(fiber.StatusNotFound, "Task not found!"))
		}
		panic(fmt.Errorf("task_repo_pg_error: delete task: %v", err))
	}
}

func (r *TaskRepositoryPG) GetTasksByOwner(ownerId string) []entity.PreviewTask {
	var tasks []entity.PreviewTask
	var project sql.NullString

	query := `SELECT t.id, t.title, t.description, t.priority, t.status, p.title as project 
			  FROM tasks t 
			  LEFT JOIN projects p ON t.project_id = p.id 
			  WHERE t.owner_id = $1`

	rows, err := r.db.Query(query, ownerId)
	if err != nil {
		panic(fmt.Errorf("task_repo_pg_error: get tasks by owner: %v", err))
	}
	log.Printf("Tasks Repo: %v", rows)
	defer rows.Close()

	for rows.Next() {
		var task entity.PreviewTask
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.Status, &project); err != nil {
			panic(fmt.Errorf("task_repo_pg_error: scan task: %v", err))
		}
		task.Project = project.String
		tasks = append(tasks, task)
	}

	return tasks
}

func (r *TaskRepositoryPG) GetTasksByAssignedUser(userId string) []entity.PreviewTask {
	var tasks []entity.PreviewTask

	query := `SELECT t.id, t.title, t.description, t.priority, t.status, p.title as project 
			  FROM tasks t 
			  JOIN task_assignments ta ON t.id = ta.task_id 
			  JOIN projects p ON t.project_id = p.id 
			  WHERE ta.user_id = $1`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		panic(fmt.Errorf("task_repo_pg_error: get tasks by assigned user: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		var task entity.PreviewTask
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.Project); err != nil {
			panic(fmt.Errorf("task_repo_pg_error: scan task: %v", err))
		}
		tasks = append(tasks, task)
	}

	return tasks
}
