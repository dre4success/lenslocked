package models

import (
	"database/sql"
	"fmt"
	"time"
)

const DefaultResetDuration = 1 * time.Hour

type PasswordReset struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
	// To determine how many bytes to use when generating password reset token 
	BytesPerToken int
	// The amoung of time that a PasswordReset is valid for
	Duration time.Duration
}

func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("T")
}

func (service *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("T")
}