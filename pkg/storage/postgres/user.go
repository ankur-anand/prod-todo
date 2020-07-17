package postgres

import (
	"context"
	"fmt"

	"github.com/ankur-anand/prod-todo/pkg/storage/serror"

	"github.com/ankur-anand/prod-todo/pkg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Compile-time check for ensuring TodoStorage implements model.UserRepository.
var _ pkg.UserStorage = (*UserStorage)(nil)

// UserStorage provides a User Storage implementation over a PostgreSQL database
type UserStorage struct {
	// db holds connection in a pool for optimal performance
	db *pgxpool.Pool
}

// NewUserStore returns an initialized AUTH UserStorage storage with connection pool
func NewUserStore(db *pgxpool.Pool) (UserStorage, error) {
	if db == nil {
		return UserStorage{}, fmt.Errorf("db proxy pool is nil")
	}
	return UserStorage{db: db}, nil
}

// // Find returns an UserModel associated with the ID in the DB
func (p UserStorage) Find(ctx context.Context, id uuid.UUID) (pkg.UserModel, error) {
	var user pkg.UserModel
	// pgx close the row for reuse
	err := p.db.QueryRow(ctx, findUserByIDQuery, id).Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Username)
	switch err {
	case nil:
		return user, nil
	case pgx.ErrNoRows:
		return user, serror.NewQueryError(findUserByIDQuery, serror.ErrUserNotFound, err.Error())
	default:
		// something else went wrong,
		return user, serror.NewQueryError(findUserByIDQuery, err, err.Error())
	}
}

// FindByEmail returns an UserModel associated with the emailID in the DB
func (p UserStorage) FindByEmail(ctx context.Context, email string) (pkg.UserModel, error) {
	var user pkg.UserModel
	// pgx close the row for reuse
	err := p.db.QueryRow(ctx, findUserByEmailQuery, email).Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Username)
	switch err {
	case nil:
		return user, nil
	case pgx.ErrNoRows:
		return user, serror.NewQueryError(findUserByEmailQuery, serror.ErrUserNotFound, err.Error())
	default:
		// something else went wrong
		return user, serror.NewQueryError(findUserByEmailQuery, err, err.Error())
	}
}

// Update stores the updated user mode inside the DB
func (p UserStorage) Update(ctx context.Context, user pkg.UserModel) error {
	cmd, err := p.db.Exec(ctx, updateUserQuery, user.ID, user.Email, user.Password, user.FirstName, user.LastName, user.Username)
	if err != nil {
		return serror.NewQueryError(updateUserQuery, err, err.Error())
	}
	if !cmd.Update() && cmd.RowsAffected() != 1 {
		return serror.NewQueryError(updateUserQuery, serror.ErrUpdateCommand, "")
	}
	return nil
}

// Store stores the user mode inside the DB
func (p UserStorage) Store(ctx context.Context, user pkg.UserModel) (uuid.UUID, error) {
	cmd, err := p.db.Exec(ctx, storeUserQuery, user.ID, user.Email, user.Password, user.FirstName, user.LastName, user.Username)
	if err != nil {
		return uuid.Nil, serror.NewQueryError(storeUserQuery, err, err.Error())
	}
	if !cmd.Insert() && cmd.RowsAffected() != 1 {
		return uuid.Nil, serror.NewQueryError(storeUserQuery, serror.ErrInsertCommand, "")
	}
	return user.ID, nil
}
