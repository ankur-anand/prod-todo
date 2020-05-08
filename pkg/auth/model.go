package auth

import "github.com/google/uuid"

// UserModel represents individual User
type UserModel struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Email     string    `json:"email"`
	Password  string    `json:"-,"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
}
