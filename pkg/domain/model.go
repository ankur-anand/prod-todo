package domain

import "github.com/google/uuid"

var (
	// NilUserModel is empty UserModel, all zeros
	NilUserModel UserModel
	// NilTodoModel is empty TodoModel, all zeros
	NilTodoModel TodoModel
)

// UserModel represents individual user registered in the system
type UserModel struct {
	ID        uuid.UUID
	Email     string
	Password  string
	FirstName string
	LastName  string
	Username  string
}

// TodoModel is each single individual task
type TodoModel struct {
	Id       uuid.UUID
	UserID   uuid.UUID
	Title    string
	Content  string
	Finished bool
}
