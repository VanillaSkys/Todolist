package router

import (
	v1 "github.com/VanillaSkys/todo_fiber/cmd/web/router/v1"
	"github.com/gofiber/fiber/v3"
)

func SetupApiRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1.SetupV1Routes(api)
}
