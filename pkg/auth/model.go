package auth

import "github.com/google/uuid"

// UserModel represents individual User
type UserModel struct {
	ID        uuid.UUID
	Email     string
	Password  string
	FirstName string
	LastName  string
	Username  string
}
