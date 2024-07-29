package models

import (
	"github.com/google/uuid"
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
	ID                     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"ID"`
	Username               string    `gorm:"uniqueIndex;not null"                            json:"username"`
	Email                  string    `gorm:"uniqueIndex;not null"                            json:"email"`
	Password               string    `gorm:"not null"                                        json:"password"`
	UserType               UserType  `gorm:"not null"                                        json:"user_type"`
	FirstName              string    `                                                       json:"first_name"`
	LastName               string    `                                                       json:"last_name"`
	IsActive               bool      `gorm:"default:true"                                    json:"is_active"`
	CanRead                bool      `gorm:"default:true"                                    json:"can_read"`
	CanWrite               bool      `gorm:"default:false"                                   json:"can_write"`
	CanModerate            bool      `gorm:"default:false"                                   json:"can_moderate"`
	IsAdmin                bool      `gorm:"default:false"                                   json:"is_admin"`
	VerificationCode       string    `gorm:"size:6"                                          json:"-"`
	VerificationCodeSentAt time.Time `                                                       json:"-"`
	FailedLoginAttempts    int       `                                                       json:"-"`
	LastFailedLogin        time.Time `                                                       json:"-"`
	LockedUntil            time.Time `                                                       json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (n *User) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
