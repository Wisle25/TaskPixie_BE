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
	return &ProjectHandler{
		useCase: useCase,
	}
}

func (h *ProjectHandler) AddProject(c *fiber.Ctx) error {
	var payload entity.ProjectPayload
	_ = c.BodyParser(&payload)

	loggedUserId := c.Locals("userInfo").(entity.User).Id

	projectId := h.useCase.ExecuteAddProject(&payload, loggedUserId)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"data":    projectId,
		"message": "Project created successfully!",
	})
}

func (h *ProjectHandler) GetProjectById(c *fiber.Ctx) error {
	id := c.Params("id")

	project := h.useCase.ExecuteGetProjectById(id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   project,
	})
}

func (h *ProjectHandler) UpdateProjectById(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload entity.ProjectPayload
	_ = c.BodyParser(&payload)

	h.useCase.ExecuteUpdateProjectById(id, &payload)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Project updated successfully!",
	})
}

func (h *ProjectHandler) DeleteProjectById(c *fiber.Ctx) error {
	id := c.Params("id")

	h.useCase.ExecuteDeleteProjectById(id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Project deleted successfully!",
	})
}

func (h *ProjectHandler) GetProjects(c *fiber.Ctx) error {
	loggedUserId := c.Locals("userInfo").(entity.User).Id

	projects := h.useCase.ExecuteGetProjects(loggedUserId)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   projects,
	})
}

func (h *ProjectHandler) GetMembersProject(c *fiber.Ctx) error {
	id := c.Params("id")

	members := h.useCase.ExecuteGetProjectMembers(id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   members,
	})
}
