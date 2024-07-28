package main

import (
	"UserAuth/internal/database"
	"UserAuth/internal/handlers"
	"UserAuth/internal/middleware"
	"UserAuth/pkg/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"log"
	"time"
)

func init() {
	utils.LoadEnv()
	database.ConnectToDb()
	database.SyncDatabase()
}

func main() {
	// set Gin to Release mode
	gin.SetMode(gin.ReleaseMode)

	// initialize gin
	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		panic(err)
	}

	// Initialize Providers Auth
	handlers.GoogleAuth()
	handlers.GithubAuth()
	r.Use(middleware.GothProvider)

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://note-taking-dusky.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Apply general rate limiter to all routes
	r.Use(middleware.RateLimiter(60, time.Minute)) // 60 requests per minute

	r.POST(
		"/login",
		middleware.RateLimiter(20, time.Minute),
		middleware.CheckNotAuthenticated(),
		handlers.Login,
	)
	r.GET("/check-auth", handlers.CheckAuth)
	r.POST(
		"/register",
		middleware.RateLimiter(10, time.Minute),
		middleware.CheckNotAuthenticated(),
		handlers.Register,
	)
	r.POST(
		"/reset-request",
		middleware.RateLimiter(8, time.Minute),
		middleware.CheckNotAuthenticated(),
		handlers.ResetRequest,
	)

	// Add the remaining routes
	r.POST("/logout", middleware.CheckAuthenticated(), handlers.Logout)
	r.GET("/protected", middleware.CheckAuthenticated(), handlers.ProtectedRoute)
	r.POST("/reset-password", middleware.CheckNotAuthenticated(), handlers.ResetPassword)
	r.GET("/user", middleware.CheckAuthenticated(), handlers.GetCurrentUser)
	r.GET("/users", middleware.CheckAuthenticated(), handlers.GetAllUsers)
	r.GET("/users/:id", middleware.CheckAuthenticated(), handlers.GetUserById)
	r.PUT("/users/:id", middleware.CheckAuthenticated(), handlers.EditUser)
	r.DELETE("/users/:id", middleware.CheckAuthenticated(), handlers.DeleteUser)
	r.GET("/auth/:provider", middleware.CheckNotAuthenticated(), func(c *gin.Context) {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	})
	r.GET("/auth/:provider/callback", handlers.CallbackHandler)

	// Start the server
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
