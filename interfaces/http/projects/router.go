package projects

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/task-pixie/applications/use_case"
	"github.com/wisle25/task-pixie/interfaces/http/middlewares"
)

func NewProjectRouter(
	app *fiber.App,
	jwtMiddleware *middlewares.JwtMiddleware,
	useCase *use_case.ProjectUseCase,
) {
	projectHandler := NewProjectHandler(useCase)

	app.Post("/projects", jwtMiddleware.GuardJWT, projectHandler.AddProject)
	app.Get("/projects/:id", jwtMiddleware.GuardJWT, projectHandler.GetProjectById)
	app.Put("/projects/:id", jwtMiddleware.GuardJWT, projectHandler.UpdateProjectById)
	app.Delete("/projects/:id", jwtMiddleware.GuardJWT, projectHandler.DeleteProjectById)
	app.Get("/projects", jwtMiddleware.GuardJWT, projectHandler.GetProjects)
	app.Get("/projects-member/:id", projectHandler.GetMembersProject)
}
