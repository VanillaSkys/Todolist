package v1

import (
	"github.com/VanillaSkys/todo_fiber/internal/adapter/in/http"
	"github.com/VanillaSkys/todo_fiber/internal/adapter/out/postgres"
	"github.com/VanillaSkys/todo_fiber/internal/adapter/out/redis"
	"github.com/VanillaSkys/todo_fiber/internal/core/domain/dto"
	"github.com/VanillaSkys/todo_fiber/internal/core/domain/service"
	"github.com/VanillaSkys/todo_fiber/internal/infrastructure"
	"github.com/gofiber/fiber/v3"
)

func SetupTodoRoutes(router fiber.Router) {
	infrastructure.Db.AutoMigrate(dto.Todo{})
	todoRepo := postgres.NewGormTodoRepository(infrastructure.Db)
	todoCache := redis.NewRedisCache(infrastructure.RedisClient)
	todoService := service.NewTodoService(todoRepo, todoCache)
	todoHttp := http.NewHttpTodo(todoService)

	todo := router.Group("/todo")

	todo.Get("/", todoHttp.FindAll)
	todo.Post("/", todoHttp.Create)
	todo.Put("/", todoHttp.Update)
	todo.Delete("/", todoHttp.Delete)
}
