package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sipodi/backend/internal/domain"
	"github.com/sipodi/backend/internal/handler"
	"github.com/sipodi/backend/internal/middleware"
	"github.com/sipodi/backend/internal/service"
)

type Router struct {
	authHandler         *handler.AuthHandler
	userHandler         *handler.UserHandler
	schoolHandler       *handler.SchoolHandler
	talentHandler       *handler.TalentHandler
	verificationHandler *handler.VerificationHandler
	notificationHandler *handler.NotificationHandler
	uploadHandler       *handler.UploadHandler
	dashboardHandler    *handler.DashboardHandler
	exportHandler       *handler.ExportHandler
	authService         *service.AuthService
}

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	schoolHandler *handler.SchoolHandler,
	talentHandler *handler.TalentHandler,
	verificationHandler *handler.VerificationHandler,
	notificationHandler *handler.NotificationHandler,
	uploadHandler *handler.UploadHandler,
	dashboardHandler *handler.DashboardHandler,
	exportHandler *handler.ExportHandler,
	authService *service.AuthService,
) *Router {
	return &Router{
		authHandler:         authHandler,
		userHandler:         userHandler,
		schoolHandler:       schoolHandler,
		talentHandler:       talentHandler,
		verificationHandler: verificationHandler,
		notificationHandler: notificationHandler,
		uploadHandler:       uploadHandler,
		dashboardHandler:    dashboardHandler,
		exportHandler:       exportHandler,
		authService:         authService,
	}
}

func (r *Router) Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/login", r.authHandler.Login)
	auth.Post("/refresh", r.authHandler.Refresh)
	auth.Post("/logout", middleware.AuthMiddleware(r.authService), r.authHandler.Logout)
	auth.Post("/logout-all", middleware.AuthMiddleware(r.authService), r.authHandler.LogoutAll)

	// Protected routes
	protected := api.Group("", middleware.AuthMiddleware(r.authService))

	// Profile routes
	protected.Get("/me", r.userHandler.GetMe)
	protected.Patch("/me", r.userHandler.UpdateMe)
	protected.Patch("/me/password", r.userHandler.ChangePassword)

	// My talents (GTK)
	protected.Get("/me/talents", r.talentHandler.ListMyTalents)
	protected.Post("/me/talents", middleware.RoleMiddleware(domain.RoleGTK), r.talentHandler.Create)
	protected.Put("/me/talents/:id", middleware.RoleMiddleware(domain.RoleGTK), r.talentHandler.Update)
	protected.Delete("/me/talents/:id", middleware.RoleMiddleware(domain.RoleGTK), r.talentHandler.Delete)

	// My notifications
	protected.Get("/me/notifications", r.notificationHandler.List)
	protected.Get("/me/notifications/unread-count", r.notificationHandler.GetUnreadCount)
	protected.Patch("/me/notifications/:id/read", r.notificationHandler.MarkAsRead)
	protected.Patch("/me/notifications/read-all", r.notificationHandler.MarkAllAsRead)

	// Schools routes
	schools := protected.Group("/schools")
	schools.Get("/", middleware.RoleMiddleware(domain.RoleSuperAdmin), r.schoolHandler.List)
	schools.Get("/:id", r.schoolHandler.GetByID)
	schools.Post("/", middleware.RoleMiddleware(domain.RoleSuperAdmin), r.schoolHandler.Create)
	schools.Put("/:id", middleware.RoleMiddleware(domain.RoleSuperAdmin), r.schoolHandler.Update)
	schools.Delete("/:id", middleware.RoleMiddleware(domain.RoleSuperAdmin), r.schoolHandler.Delete)
	schools.Get("/:id/users", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.schoolHandler.GetUsers)

	// Users routes
	users := protected.Group("/users")
	users.Get("/", middleware.RoleMiddleware(domain.RoleSuperAdmin), r.userHandler.List)
	users.Get("/:id", r.userHandler.GetByID)
	users.Post("/", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.userHandler.Create)
	users.Put("/:id", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.userHandler.Update)
	users.Delete("/:id", middleware.RoleMiddleware(domain.RoleSuperAdmin), r.userHandler.Delete)
	users.Patch("/:id/activate", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.userHandler.Activate)
	users.Patch("/:id/deactivate", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.userHandler.Deactivate)

	// Talents routes
	talents := protected.Group("/talents")
	talents.Get("/", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.talentHandler.List)
	talents.Get("/:id", r.talentHandler.GetByID)

	// Verification routes
	verifications := protected.Group("/verifications")
	verifications.Get("/talents", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.verificationHandler.ListPending)
	// Batch routes MUST come BEFORE parameterized routes to avoid :id matching "batch"
	verifications.Post("/talents/batch/approve", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.verificationHandler.BatchApprove)
	verifications.Post("/talents/batch/reject", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.verificationHandler.BatchReject)
	verifications.Post("/talents/:id/approve", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.verificationHandler.Approve)
	verifications.Post("/talents/:id/reject", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.verificationHandler.Reject)

	// Upload routes
	uploads := protected.Group("/uploads")
	uploads.Post("/presign", r.uploadHandler.Presign)
	uploads.Post("/:upload_id/confirm", r.uploadHandler.Confirm)
	uploads.Delete("/:upload_id", r.uploadHandler.Cancel)

	// Dashboard routes
	protected.Get("/dashboard/summary", r.dashboardHandler.GetSummary)
	protected.Get("/dashboard/schools/statistics", middleware.RoleMiddleware(domain.RoleSuperAdmin), r.dashboardHandler.GetSchoolsStatistics)
	protected.Get("/dashboard/talents/statistics", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.dashboardHandler.GetTalentsStatistics)

	// Export routes
	exports := protected.Group("/exports")
	exports.Get("/gtk", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.exportHandler.ExportGTK)
	exports.Get("/talents", middleware.RoleMiddleware(domain.RoleSuperAdmin, domain.RoleAdminSekolah), r.exportHandler.ExportTalents)
	exports.Get("/schools", middleware.RoleMiddleware(domain.RoleSuperAdmin), r.exportHandler.ExportSchools)
}
