package pkg

import (
	"context"

	"github.com/google/uuid"
)

var (
	// NilTodoModel is empty TodoModel, all zeros
	NilTodoModel TodoModel
)

// TodoModel is each single individual task
type TodoModel struct {
	Title    string
	Content  string
	ID       uuid.UUID
	UserID   uuid.UUID
	Finished bool
}

// TodoFilter tells what kind of filter to apply on queries
type TodoFilter int8

const (
	// NilFilter rule suggests no filter should be applied on the query
	NilFilter TodoFilter = iota
	// Finished rule to filter all the Finished=True Todo
	Finished
	// UnFinished rule to filter all the Finished=False Todo
	UnFinished
)

// TodoStorage define a contract for storage, to interact
// with the Todo Model.
type TodoStorage interface {
	FindOneTodo(ctx context.Context, id uuid.UUID) (TodoModel, error)
	FindAllTodoOfUser(ctx context.Context, userID uuid.UUID, filter TodoFilter) ([]TodoModel, error)
	UpdateOne(ctx context.Context, todo TodoModel) error
	InsertOne(ctx context.Context, todo TodoModel) (uuid.UUID, error)
	DeleteOne(ctx context.Context, id uuid.UUID) error
}
