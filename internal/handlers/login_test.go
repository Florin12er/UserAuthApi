package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})

	return db
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/login", Login)
	return r
}

func TestLogin(t *testing.T) {
	// Setup
	database.DB = setupTestDB()
	router := setupRouter()

	// Create a test user with the specified credentials
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := models.User{
		Username: "florin",
		Email:    "florin@example.com",
		Password: string(hashedPassword),
	}
	database.DB.Create(&testUser)

	// Test cases
	tests := []struct {
		name           string
		payload        gin.H
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid Login",
			payload: gin.H{
				"email_or_username": "florin",
				"password":          "password123",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Logged in successfully",
		},
		{
			name: "Invalid Password",
			payload: gin.H{
				"email_or_username": "florin",
				"password":          "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid credentials",
		},
		{
			name: "Non-existent User",
			payload: gin.H{
				"email_or_username": "nonexistent",
				"password":          "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

