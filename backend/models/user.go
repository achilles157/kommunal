package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"` // "-" means this field won't be included in JSON responses
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// SignUpInput represents the data needed for user registration
type SignUpInput struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Username string `json:"username" binding:"required,min=3,max=30" example:"johndoe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// SignInInput represents the data needed for user authentication
type SignInInput struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// UpdateProfileInput represents the data that can be updated in a user's profile
type UpdateProfileInput struct {
	Name     string `json:"name" binding:"omitempty,min=2"`
	Username string `json:"username" binding:"omitempty,min=3,max=30"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}
