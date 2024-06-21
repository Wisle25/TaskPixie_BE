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
	var task entity.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request payload",
		})
	}
	task.UserId = c.Locals("userInfo").(entity.User).Id

	taskId := h.useCase.ExecuteAddTask(&task)

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

func (h *TaskHandler) GetTasksByUser(c *fiber.Ctx) error {
	userId := c.Params("userId")
	tasks := h.useCase.ExecuteGetTasksByUser(userId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   tasks,
	})
}

func (h *TaskHandler) GetTasksByProject(c *fiber.Ctx) error {
	projectId := c.Params("projectId")
	tasks := h.useCase.ExecuteGetTasksByProject(projectId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   tasks,
	})
}

func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	var task entity.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request payload",
		})
	}

	taskId := h.useCase.ExecuteUpdateTask(&task)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"data":    taskId,
		"message": "Task updated successfully",
	})
}

func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	h.useCase.ExecuteDeleteTask(id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Task deleted successfully",
	})
}
