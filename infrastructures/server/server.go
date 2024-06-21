package server

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/wisle25/task-pixie/commons"
	"github.com/wisle25/task-pixie/infrastructures/cache"
	"github.com/wisle25/task-pixie/infrastructures/container"
	"github.com/wisle25/task-pixie/infrastructures/file_statics"
	"github.com/wisle25/task-pixie/infrastructures/generator"
	"github.com/wisle25/task-pixie/infrastructures/services"
	"github.com/wisle25/task-pixie/interfaces/http/middlewares"
	"github.com/wisle25/task-pixie/interfaces/http/projects"
	"github.com/wisle25/task-pixie/interfaces/http/tasks"
	"github.com/wisle25/task-pixie/interfaces/http/users"
)

func errorHandling(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	status := "error"
	code := fiber.StatusInternalServerError
	message := err.Error()

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		status = "fail"
		code = e.Code
		message = e.Message
	}

	// Send custom error
	return c.Status(code).JSON(fiber.Map{
		"status":  status,
		"message": message,
	})
}

func CreateServer(config *commons.Config) *fiber.App {
	// Load Services
	db := services.ConnectDB(config)
	redis := services.ConnectRedis(config)
	minio, bucketName := services.NewMinio(config)

	// Server
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandling,
	})

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "POST,GET,PUT,DELETE",
	}))

	// Global Dependencies
	redisCache := cache.NewRedisCache(redis)
	uuidGenerator := generator.NewUUIDGenerator()
	validation := services.NewValidation()
	minioFileUpload := file_statics.NewMinioFileUpload(minio, uuidGenerator, bucketName)
	vipsFileProcessing := file_statics.NewVipsFileProcessing()

	// Use Cases
	userUseCase := container.NewUserContainer(
		config,
		db,
		redisCache,
		uuidGenerator,
		vipsFileProcessing,
		minioFileUpload,
		validation,
	)
	projectUseCase := container.NewProjectContainer(uuidGenerator, db)
	tasksUseCase := container.NewTaskContainer(uuidGenerator, db)

	// Custom Middleware
	jwtMiddleware := middlewares.NewJwtMiddleware(userUseCase)

	// Router
	users.NewUserRouter(app, jwtMiddleware, userUseCase)
	projects.NewProjectRouter(app, jwtMiddleware, projectUseCase)
	tasks.NewTaskRouter(app, jwtMiddleware, tasksUseCase)

	return app
}
