package auth

import (
	"context"
)

// Repository define a repository, where the entities will be stored.
type Repository interface {
	Find(ctx context.Context, id int) (UserModel, error)
	FindByEmail(ctx context.Context, email string) (UserModel, error)
	FindAll(ctx context.Context) ([]UserModel, error)
	Update(ctx context.Context, user UserModel) error
	Store(ctx context.Context, user UserModel) (int, error)
}
