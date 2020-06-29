package model

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository define a contract for storage, to interact
// with the UserModel.
type UserRepository interface {
	Find(ctx context.Context, id uuid.UUID) (UserModel, error)
	FindByEmail(ctx context.Context, email string) (UserModel, error)
	FindAll(ctx context.Context) (UserIterator, error)
	Update(ctx context.Context, user UserModel) error
	Store(ctx context.Context, user UserModel) (uuid.UUID, error)
}

// Iterator is implemented by type that can be iterated.
// As there is no upper bound in the number of users, we
// fetch lazily on demand.
type Iterator interface {
	// Next advances the iterator. If no more items are available or an
	// error occurs, calls to Next() return false.
	Next() bool

	// Error returns the last error encountered by the iterator.
	Error() error

	// Close releases any resources associated with an iterator.
	Close() error
}

// UserIterator is implemented by type that can iterate the Users.
type UserIterator interface {
	Iterator

	// User returns the currently fetched User Model.
	User() UserModel
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

// TodoRepository define a contract for storage, to interact
// with the Todo Model.
type TodoRepository interface {
	FindOneTodo(ctx context.Context, id uuid.UUID) (TodoModel, error)
	FindAllTodoOfUser(ctx context.Context, userID uuid.UUID, filter TodoFilter) (TodoIterator, error)
	FindAllByFilter(ctx context.Context, filter TodoFilter) (TodoIterator, error)
	UpdateOne(ctx context.Context, todo TodoModel) error
	UpdateMany(ctx context.Context, todo []*TodoModel) error
	InsertOne(ctx context.Context, todo TodoModel) (uuid.UUID, error)
	InsertMany(ctx context.Context, todo []*TodoModel) (uuid.UUID, error)
	DeleteOne(ctx context.Context, id uuid.UUID) error
	DeleteMany(ctx context.Context, id []uuid.UUID) error
}

// TodoIterator is implemented by type that can iterate the Todo List.
type TodoIterator interface {
	Iterator

	// Todo returns the currently fetched Todo Model.
	Todo() TodoModel
}
