package main

import (
	"UserAuth/internal/database"
	"UserAuth/internal/handlers"
	"UserAuth/internal/middleware"
	"UserAuth/pkg/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	utils.LoadEnv()
	database.ConnectToDb()
	database.SyncDatabase()
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	// initialize gin
	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		panic(err)
	}

	// Apply general rate limiter to all routes
	r.Use(middleware.RateLimiter(60, time.Minute)) // 60 requests per minute

	// Apply specific rate limiters to sensitive endpoints
	r.POST("/login", middleware.RateLimiter(5, time.Minute), middleware.CheckNotAuthenticated(), handlers.Login)
	r.POST("/register", middleware.RateLimiter(3, time.Minute), middleware.CheckNotAuthenticated(), handlers.Register)
	r.POST("/reset-request", middleware.RateLimiter(2, time.Minute), middleware.CheckNotAuthenticated(), handlers.ResetRequest)

	// Add the remaining routes
	r.POST("/logout", middleware.CheckAuthenticated(), handlers.Logout)
	r.POST("/verify-email", handlers.VerifyEmail)
	r.GET("/protected", middleware.CheckAuthenticated(), handlers.ProtectedRoute)
	r.POST("/reset-password", middleware.CheckNotAuthenticated(), handlers.ResetPassword)
	r.GET("/user", middleware.CheckAuthenticated(), handlers.GetCurrentUser)
	r.GET("/users", middleware.CheckAuthenticated(), handlers.GetAllUsers)
	r.GET("/users/:id", middleware.CheckAuthenticated(), handlers.GetUserById)
	r.PUT("/users/:id", middleware.CheckAuthenticated(), handlers.EditUser)
	r.DELETE("/users/:id", middleware.CheckAuthenticated(), handlers.DeleteUser)

	// Start the server
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

