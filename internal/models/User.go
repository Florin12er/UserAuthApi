package models

import (
	"gorm.io/gorm"
	"time"
)

type UserType string

const (
	TypeMember    UserType = "member"
	TypeGuest     UserType = "guest"
	TypeAdmin     UserType = "admin"
	TypeModerator UserType = "moderator"
)

type User struct {
	gorm.Model
	Username               string    `json:"username"     gorm:"uniqueIndex;not null"`
	Email                  string    `json:"email"        gorm:"uniqueIndex;not null"`
	Password               string    `json:"password"     gorm:"not null"`
	UserType               UserType  `json:"user_type"    gorm:"not null"`
	FirstName              string    `json:"first_name"`
	LastName               string    `json:"last_name"`
	IsActive               bool      `json:"is_active"    gorm:"default:true"`
	CanRead                bool      `json:"can_read"     gorm:"default:true"`
	CanWrite               bool      `json:"can_write"    gorm:"default:false"`
	CanModerate            bool      `json:"can_moderate" gorm:"default:false"`
	IsAdmin                bool      `json:"is_admin"     gorm:"default:false"`
	VerificationCode       string    `json:"-"            gorm:"size:6"`
	VerificationCodeSentAt time.Time `json:"-"`
	FailedLoginAttempts    int       `json:"-"`
	LastFailedLogin        time.Time `json:"-"`
	LockedUntil            time.Time `json:"-"`
}
