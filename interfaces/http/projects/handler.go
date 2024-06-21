package projects

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/task-pixie/applications/use_case"
	"github.com/wisle25/task-pixie/domains/entity"
)

type ProjectHandler struct {
	useCase *use_case.ProjectUseCase
}

func NewProjectHandler(useCase *use_case.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{useCase: useCase}
}

func (h *ProjectHandler) AddProject(c *fiber.Ctx) error {
	// Payload
	var payload entity.ProjectPayload
	_ = c.BodyParser(&payload)

	ownerId := c.Locals("userInfo").(entity.User).Id

	// Use Case
	returnedId := h.useCase.ExecuteAddProject(&payload, ownerId)

	// Response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"data":    returnedId,
		"message": "Successfully created new project!",
	})
}

func (h *ProjectHandler) GetProjectById(c *fiber.Ctx) error {
	// Payload
	id := c.Params("id")

	// Use Case
	project := h.useCase.ExecuteGetProjectById(id)

	// Response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   project,
	})
}

func (h *ProjectHandler) UpdateProjectById(c *fiber.Ctx) error {
	// Payload
	id := c.Params("id")
	var payload entity.ProjectPayload
	_ = c.BodyParser(&payload)

	// Use Case
	h.useCase.ExecuteUpdateProjectById(id, &payload)

	// Response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully updated project!",
	})
}

func (h *ProjectHandler) DeleteProjectById(c *fiber.Ctx) error {
	// Payload
	id := c.Params("id")

	// Use Case
	h.useCase.ExecuteDeleteProjectById(id)

	// Response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully deleted project!",
	})
}

func (h *ProjectHandler) GetProjectsByOwner(c *fiber.Ctx) error {
	ownerId := c.Locals("userInfo").(entity.User).Id

	// Use Case
	projects := h.useCase.ExecuteGetProjectsByOwner(ownerId)

	// Response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   projects,
	})
}
