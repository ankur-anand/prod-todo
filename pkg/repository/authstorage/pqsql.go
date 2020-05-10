package authstorage

import (
	"context"
	"fmt"

	"github.com/ankur-anand/prod-app/pkg/repository/storageerror"

	"github.com/ankur-anand/prod-app/pkg/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	// SQL Query
	findUserByIDQuery    = "SELECT * FROM users WHERE user_id=$1"
	findUserByEmailQuery = "SELECT * FROM users WHERE email_id=$1"

	storeUserQuery = `
INSERT INTO users (user_id, email_id, password_hash, first_name, last_name, user_name) VALUES ($1, $2, $3, $4, $5, $6)
`
)

// PgSQL provides a AUTH Repository implementation over a PostgreSQL database
type PgSQL struct {
	// db holds connection in a pool for optimal performance
	db *pgxpool.Pool
}

// NewAuthPgSQL returns an initialized AUTH PgSQL repository with connection pool
func NewAuthPgSQL(db *pgxpool.Pool) (PgSQL, error) {
	if db == nil {
		return PgSQL{}, fmt.Errorf("db proxy pool is nil")
	}
	return PgSQL{db: db}, nil
}

// Find returns an UserModel associated with the ID in the DB
func (p PgSQL) Find(ctx context.Context, id uuid.UUID) (auth.UserModel, error) {
	var user auth.UserModel
	// pgx close the row for reuse
	err := p.db.QueryRow(ctx, findUserByIDQuery, id).Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Username)
	switch err {
	case nil:
		return user, nil
	case pgx.ErrNoRows:
		return user, storageerror.NewQueryError(findUserByIDQuery, ErrUserNotFound, err.Error())
	default:
		// something else went wrong,
		return user, storageerror.NewQueryError(findUserByIDQuery, err, err.Error())
	}
}

// FindByEmail returns an UserModel associated with the emailID in the DB
func (p PgSQL) FindByEmail(ctx context.Context, email string) (auth.UserModel, error) {
	var user auth.UserModel
	// pgx close the row for reuse
	err := p.db.QueryRow(ctx, findUserByEmailQuery, email).Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Username)
	switch err {
	case nil:
		return user, nil
	case pgx.ErrNoRows:
		return user, storageerror.NewQueryError(findUserByEmailQuery, ErrUserNotFound, err.Error())
	default:
		// something else went wrong
		return user, storageerror.NewQueryError(findUserByEmailQuery, err, err.Error())
	}
}

// FindAll returns all User inside the DB
func (p PgSQL) FindAll(ctx context.Context) (auth.UserIterator, error) {
	panic("implement me")
}

// Update stores the updated user mode inside the DB
func (p PgSQL) Update(ctx context.Context, user auth.UserModel) error {
	panic("implement me")
}

// Store stores the user mode inside the DB
func (p PgSQL) Store(ctx context.Context, user auth.UserModel) (uuid.UUID, error) {
	cmd, err := p.db.Exec(ctx, storeUserQuery, user.ID, user.Email, user.Password, user.FirstName, user.LastName, user.Username)
	if err != nil {
		return uuid.Nil, storageerror.NewQueryError(storeUserQuery, err, err.Error())
	}
	if !cmd.Insert() && cmd.RowsAffected() != 1 {
		return uuid.Nil, storageerror.NewQueryError(storeUserQuery, ErrInsertCommand, "")
	}
	return user.ID, nil
}
