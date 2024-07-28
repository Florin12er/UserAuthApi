package handlers

import (
	"UserAuth/internal/database"
	"UserAuth/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForRegister() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})

	return db
}

func setupRouterForRegister() *gin.Engine {
	r := gin.Default()
	r.POST("/register", Register)
	return r
}

func TestRegister(t *testing.T) {
	// Setup
	database.DB = setupTestDBForRegister()
	router := setupRouterForRegister()

	// Set test environment
	os.Setenv("APP_ENV", "test")

	tests := []struct {
		name           string
		payload        gin.H
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid Registration",
			payload: gin.H{
				"username":              "testuser",
				"email":                 "test@example.com",
				"password":              "Password123!",
				"cf-turnstile-response": "mock_turnstile_response",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Registration successful",
		},
		{
			name: "Duplicate Username",
			payload: gin.H{
				"username":              "testuser",
				"email":                 "another@example.com",
				"password":              "Password123!",
				"cf-turnstile-response": "mock_turnstile_response",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Username or email already exists",
		},
		{
			name: "Invalid Email",
			payload: gin.H{
				"username":              "newuser",
				"email":                 "invalid-email",
				"password":              "Password123!",
				"cf-turnstile-response": "mock_turnstile_response",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Key: 'email' Error:Field validation for 'email' failed",
		},
		{
			name: "Weak Password",
			payload: gin.H{
				"username":              "newuser",
				"email":                 "new@example.com",
				"password":              "weak",
				"cf-turnstile-response": "mock_turnstile_response",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Password does not meet the security requirements",
		},
		{
			name: "Invalid CAPTCHA",
			payload: gin.H{
				"username":              "newuser",
				"email":                 "new@example.com",
				"password":              "Password123!",
				"cf-turnstile-response": "invalid_response",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid CAPTCHA",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"Valid Password", "Password123!", true},
		{"No Uppercase", "password123!", false},
		{"No Lowercase", "PASSWORD123!", false},
		{"No Number", "Password!", false},
		{"No Special", "Password123", false},
		{"Too Short", "Pass1!", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidPassword(tt.password)
			assert.Equal(t, tt.want, got)
		})
	}
}

