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

type TaskRepositoryPG struct {
	idGenerator generator.IdGenerator
	db          *sql.DB
}

func NewTaskRepositoryPG(idGenerator generator.IdGenerator, db *sql.DB) repository.TaskRepository {
	return &TaskRepositoryPG{idGenerator, db}
}

func (r *TaskRepositoryPG) AddTask(task *entity.Task) string {
	id := r.idGenerator.Generate()

	query := `
		INSERT INTO tasks (id, title, description, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var returnedId string
	err := r.db.QueryRow(query, id, task.Title, task.Description, task.UserId).Scan(&returnedId)
	if err != nil {
		panic(fmt.Errorf("task_repo_pg_error: add task: %v", err))
	}

	return returnedId
}

func (r *TaskRepositoryPG) GetTaskById(id string) *entity.Task {
	var task entity.Task
	query := `
		SELECT id, title, description, project_id, user_id, completed
		FROM tasks
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.ProjectId,
		&task.UserId,
		&task.Completed,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			panic(fiber.NewError(fiber.StatusNotFound, "Task not found!"))
		}
		panic(fmt.Errorf("task_repo_pg_error: get task by id: %v", err))
	}

	return &task
}

func (r *TaskRepositoryPG) GetTasksByUser(userId string) []entity.Task {
	query := `
		SELECT id, title, description, user_id, completed
		FROM tasks
		WHERE user_id = $1`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		panic(fmt.Errorf("task_repo_pg_error: get tasks by user: %v", err))
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.UserId,
			&task.Completed,
		)
		if err != nil {
			panic(fmt.Errorf("task_repo_pg_error: scan task: %v", err))
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func (r *TaskRepositoryPG) GetTasksByProject(projectId string) []entity.Task {
	query := `
		SELECT id, title, description, project_id, user_id, completed
		FROM tasks
		WHERE project_id = $1`

	rows, err := r.db.Query(query, projectId)
	if err != nil {
		panic(fmt.Errorf("task_repo_pg_error: get tasks by project: %v", err))
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.ProjectId,
			&task.UserId,
			&task.Completed,
		)
		if err != nil {
			panic(fmt.Errorf("task_repo_pg_error: scan task: %v", err))
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func (r *TaskRepositoryPG) UpdateTask(task *entity.Task) string {
	query := `
		UPDATE tasks
		SET title = $2, description = $3, project_id = $4, completed = $5
		WHERE id = $1
		RETURNING id`

	var updatedId string
	err := r.db.QueryRow(query, task.Id, task.Title, task.Description, task.ProjectId, task.Completed).Scan(&updatedId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			panic(fiber.NewError(fiber.StatusNotFound, "Task not found!"))
		}
		panic(fmt.Errorf("task_repo_pg_error: update task: %v", err))
	}

	return updatedId
}

func (r *TaskRepositoryPG) DeleteTask(id string) {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		panic(fmt.Errorf("task_repo_pg_error: delete task: %v", err))
	}
}
