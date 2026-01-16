package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sipodi/backend/internal/config"
	"github.com/sipodi/backend/internal/database"
	"github.com/sipodi/backend/internal/handler"
	"github.com/sipodi/backend/internal/middleware"
	"github.com/sipodi/backend/internal/repository"
	"github.com/sipodi/backend/internal/router"
	"github.com/sipodi/backend/internal/service"
	"github.com/sipodi/backend/internal/storage"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect to database
	db, err := database.NewPostgresPool(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// Connect to MinIO
	minioStorage, err := storage.NewMinIOStorage(cfg.MinIO)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}
	log.Println("Connected to MinIO")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	schoolRepo := repository.NewSchoolRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	talentRepo := repository.NewTalentRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, tokenRepo, cfg.JWT)
	userService := service.NewUserService(userRepo, schoolRepo)
	schoolService := service.NewSchoolService(schoolRepo, userRepo)
	talentService := service.NewTalentService(talentRepo, userRepo, notificationRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	uploadService := service.NewUploadService(minioStorage)
	dashboardService := service.NewDashboardService(userRepo, schoolRepo, talentRepo, notificationRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	schoolHandler := handler.NewSchoolHandler(schoolService)
	talentHandler := handler.NewTalentHandler(talentService, uploadService)
	verificationHandler := handler.NewVerificationHandler(talentService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	uploadHandler := handler.NewUploadHandler(uploadService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	exportHandler := handler.NewExportHandler(userService, schoolService, talentService)

	// Initialize router
	r := router.NewRouter(
		authHandler,
		userHandler,
		schoolHandler,
		talentHandler,
		verificationHandler,
		notificationHandler,
		uploadHandler,
		dashboardHandler,
		exportHandler,
		authService,
	)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: errorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.CORSMiddleware(cfg.CORS.Origins))

	// Setup routes
	r.Setup(app)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		app.Shutdown()
	}()

	// Start server
	log.Printf("Server starting on port %s", cfg.App.Port)
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    "INTERNAL_ERROR",
			"message": "Terjadi kesalahan pada server",
		},
	})
}
