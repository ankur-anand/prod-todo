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

type TodoList struct {
	// ID of the List
	ID uuid.UUID
	// UserID is fk to user
	UserID uuid.UUID
	// TodoID is fk to Todo
	TodoID uuid.UUID
}

type Todo struct {
	Id       int
	Title    string
	Content  string
	Finished bool
}
