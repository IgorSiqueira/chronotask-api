package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware_AllowedOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	allowedOrigins := []string{
		"https://chronotask.wizardtech.com.br",
		"http://localhost:5173",
	}
	middleware := NewCORSMiddleware(allowedOrigins)

	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// Test allowed origin
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://chronotask.wizardtech.com.br")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "https://chronotask.wizardtech.com.br" {
		t.Errorf("Expected Access-Control-Allow-Origin header, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Errorf("Expected Access-Control-Allow-Credentials: true")
	}

	if w.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
		t.Errorf("Expected Access-Control-Allow-Methods header")
	}

	if w.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Errorf("Expected Access-Control-Allow-Headers header")
	}
}

func TestCORSMiddleware_DisallowedOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	allowedOrigins := []string{
		"https://chronotask.wizardtech.com.br",
	}
	middleware := NewCORSMiddleware(allowedOrigins)

	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// Test disallowed origin
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should still return 200 but without CORS headers
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("Should not set CORS headers for disallowed origin")
	}
}

func TestCORSMiddleware_PreflightRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	allowedOrigins := []string{
		"https://chronotask.wizardtech.com.br",
	}
	middleware := NewCORSMiddleware(allowedOrigins)

	router := gin.New()
	router.Use(middleware.Handler())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// Test OPTIONS preflight request
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://chronotask.wizardtech.com.br")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != 204 {
		t.Errorf("Expected status 204 for OPTIONS, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "https://chronotask.wizardtech.com.br" {
		t.Errorf("Expected Access-Control-Allow-Origin header for preflight")
	}

	if w.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
		t.Errorf("Expected Access-Control-Allow-Methods header for preflight")
	}
}

func TestCORSMiddleware_LocalhostOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	allowedOrigins := []string{
		"http://localhost:5173",
	}
	middleware := NewCORSMiddleware(allowedOrigins)

	router := gin.New()
	router.Use(middleware.Handler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// Test localhost origin (development)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Errorf("Expected localhost origin to be allowed")
	}
}
