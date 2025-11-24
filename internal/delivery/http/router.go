package http

import (
	"github.com/gin-gonic/gin"
)

// Router holds all handlers and configures routes
type Router struct {
	healthHandler *HealthHandler
	userHandler   *UserHandler
}

// NewRouter creates a new Router with all handlers
func NewRouter(healthHandler *HealthHandler, userHandler *UserHandler) *Router {
	return &Router{
		healthHandler: healthHandler,
		userHandler:   userHandler,
	}
}

// SetupRoutes configures all application routes
func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", r.healthHandler.Check)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// User routes
		v1.POST("/user", r.userHandler.Create)
	}

	return router
}
