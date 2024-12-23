package v1

import "github.com/gofiber/fiber/v3"

func SetupV1Routes(app fiber.Router) {
	v1 := app.Group("/v1")

	SetupTodoRoutes(v1)
}
