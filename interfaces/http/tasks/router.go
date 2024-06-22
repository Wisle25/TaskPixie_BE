package tasks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/task-pixie/applications/use_case"
	"github.com/wisle25/task-pixie/interfaces/http/middlewares"
)

func NewTaskRouter(
	app *fiber.App,
	jwtMiddleware *middlewares.JwtMiddleware,
	useCase *use_case.TaskUseCase,
) {
	taskHandler := NewTaskHandler(useCase)

	app.Post("/tasks", jwtMiddleware.GuardJWT, taskHandler.AddTask)
	app.Get("/tasks/:id", jwtMiddleware.GuardJWT, taskHandler.GetTaskById)
	app.Get("/tasks", jwtMiddleware.GuardJWT, taskHandler.GetTasks)
	app.Get("/tasks/project/:projectId", jwtMiddleware.GuardJWT, taskHandler.GetTasksByProject)
	app.Put("/tasks/:id", jwtMiddleware.GuardJWT, taskHandler.UpdateTask)
	app.Delete("/tasks/:id", jwtMiddleware.GuardJWT, taskHandler.DeleteTask)
}
