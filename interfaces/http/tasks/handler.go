package tasks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wisle25/task-pixie/applications/use_case"
	"github.com/wisle25/task-pixie/domains/entity"
)

type TaskHandler struct {
	useCase *use_case.TaskUseCase
}

func NewTaskHandler(useCase *use_case.TaskUseCase) *TaskHandler {
	return &TaskHandler{useCase: useCase}
}

func (h *TaskHandler) AddTask(c *fiber.Ctx) error {
	var task entity.TaskPayload
	_ = c.BodyParser(&task)

	userId := c.Locals("userInfo").(entity.User).Id

	taskId := h.useCase.ExecuteAddTask(&task, userId)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"data":    taskId,
		"message": "Task created successfully",
	})
}

func (h *TaskHandler) GetTaskById(c *fiber.Ctx) error {
	id := c.Params("id")
	task := h.useCase.ExecuteGetTaskById(id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   task,
	})
}

func (h *TaskHandler) GetTasks(c *fiber.Ctx) error {
	userId := c.Locals("userInfo").(entity.User).Id

	tasks := h.useCase.ExecuteGetTasks(userId)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   tasks,
	})
}

func (h *TaskHandler) GetTasksByProject(c *fiber.Ctx) error {
	projectId := c.Params("projectId")

	tasks := h.useCase.ExecuteGetTasksByProjects(projectId)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   tasks,
	})
}

func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")

	var payload entity.TaskPayload
	_ = c.BodyParser(&payload)

	h.useCase.ExecuteUpdateTaskById(id, &payload)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Task updated successfully",
	})
}

func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	h.useCase.ExecuteDeleteTaskById(id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Task deleted successfully",
	})
}
