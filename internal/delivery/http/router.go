package http

import (
	"github.com/gin-gonic/gin"
	"github.com/igor/chronotask-api/internal/delivery/http/middleware"
)

// Router holds all handlers and configures routes
type Router struct {
	healthHandler  *HealthHandler
	userHandler    *UserHandler
	authMiddleware *middleware.AuthMiddleware
}

// NewRouter creates a new Router with all handlers
func NewRouter(
	healthHandler *HealthHandler,
	userHandler *UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		healthHandler:  healthHandler,
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
	}
}

// SetupRoutes configures all application routes
func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Health check endpoint (public)
	router.GET("/health", r.healthHandler.Check)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		v1.POST("/user", r.userHandler.Create)    // Create user (register)
		v1.POST("/login", r.userHandler.Login)    // Login

		// Protected routes (authentication required)
		authenticated := v1.Group("")
		authenticated.Use(r.authMiddleware.RequireAuth())
		{
			// User protected routes
			authenticated.GET("/user/profile", r.userHandler.GetProfile)

			// Future protected routes
			// authenticated.PUT("/user/profile", r.userHandler.UpdateProfile)
			// authenticated.POST("/habit", r.habitHandler.Create)
		}
	}

	return router
}
