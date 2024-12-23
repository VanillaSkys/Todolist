package main

import (
	"log"
	"strings"
	"time"

	"github.com/VanillaSkys/todo_fiber/cmd/web/router"
	"github.com/VanillaSkys/todo_fiber/internal/infrastructure"
	"github.com/VanillaSkys/todo_fiber/internal/logger"
	"github.com/VanillaSkys/todo_fiber/internal/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/spf13/viper"
)

func main() {
	app := fiber.New()

	if err := InitConfig(); err != nil {
		log.Fatalf("Error config: %v", err)
	}

	if err := logger.InitLogger(); err != nil {
		log.Fatalf("Error logger: %v", err)
	}
	defer logger.SyncLogger()

	InitTimeZone()
	infrastructure.InitPostgres()
	infrastructure.InitRedis()

	app.Use(cors.New())
	app.Use(middleware.SetRequestId())

	router.SetupApiRoutes(app)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func InitTimeZone() {
	timeZone := viper.GetString("system.timezone")
	ict, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Fatalf("Error loading timezone: %v", err)
	}
	time.Local = ict
}

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.EnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
