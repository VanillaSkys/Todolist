package http

import (
	"github.com/VanillaSkys/todo_fiber/internal/core/domain/dto"
	"github.com/VanillaSkys/todo_fiber/internal/core/domain/service"
	"github.com/VanillaSkys/todo_fiber/internal/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type httpTodoImpl struct {
	service   service.TodoService
	validator *validator.Validate
}

func NewHttpTodo(service service.TodoService) *httpTodoImpl {
	return &httpTodoImpl{service: service, validator: validator.New()}
}

func (h *httpTodoImpl) FindAll(c fiber.Ctx) error {
	requestId := c.Locals("X-Request-ID").(string)
	httpLogger := logger.Log.With(zap.String("X-Request-ID", requestId))

	httpLogger.Info("Call interface to find all todos.")
	todos, err := h.service.FindAll()
	if err != nil {
		httpLogger.Error("Error fetching todos from service", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch todos",
		})
	}

	httpLogger.Info("Returning todos.")
	return c.JSON(fiber.Map{"message": todos, "X-Request-ID": requestId})
}

func (h *httpTodoImpl) Create(c fiber.Ctx) error {
	requestId := c.Locals("X-Request-ID").(string)
	httpLogger := logger.Log.With(zap.String("X-Request-ID", requestId))

	httpLogger.Info("Call interface to create todo.")
	var input dto.TodoInputSave
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body."})
	}
	if err := h.validator.Struct(input); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return c.Status(400).JSON(fiber.Map{"error": "Validation error on field: " + e.StructField() + " - " + e.Tag()})
		}
	}
	todo := dto.Todo{
		Id:          uuid.NewString(),
		Topic:       input.Topic,
		Description: input.Description,
		Status:      input.Status,
	}

	if err := h.service.Create(todo); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error."})
	}

	httpLogger.Info("Todo created successfully.")
	return c.JSON(fiber.Map{
		"message":   "insert ok",
		"dataAdded": todo,
	})
}

func (h *httpTodoImpl) Update(c fiber.Ctx) error {
	requestId := c.Locals("X-Request-ID").(string)
	httpLogger := logger.Log.With(zap.String("X-Request-ID", requestId))

	httpLogger.Info("Call interface to update todo.")
	var input dto.TodoInputUpdateStatus
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body."})
	}
	if err := h.validator.Struct(input); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return c.Status(400).JSON(fiber.Map{"error": "Validation error on field:" + e.StructField() + " - " + e.Tag()})
		}
	}
	if err := h.service.Update(input); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error."})
	}

	httpLogger.Info("Todo updated successfully.")
	return c.JSON(fiber.Map{
		"message": "update ok",
	})
}

func (h *httpTodoImpl) Delete(c fiber.Ctx) error {
	requestId := c.Locals("X-Request-ID").(string)
	httpLogger := logger.Log.With(zap.String("X-Request-ID", requestId))

	httpLogger.Info("Call interface to delete todo.")
	var input dto.TodoInputDelete
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body."})
	}
	if err := h.validator.Struct(input); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return c.Status(400).JSON(fiber.Map{"error": "Validation error on field:" + e.StructField() + " - " + e.Tag()})
		}
	}
	if err := h.service.Delete(input); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error."})
	}
	httpLogger.Info("Todo deleted successfully.")
	return c.JSON(fiber.Map{
		"message": "deleted ok",
	})
}
