package domain

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
