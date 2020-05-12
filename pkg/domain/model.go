package domain

import "github.com/google/uuid"

// UserModel represents individual user registered in the system
type UserModel struct {
	ID        uuid.UUID
	Email     string
	Password  string
	FirstName string
	LastName  string
	Username  string
}
