package storage

import (
	"context"

	"github.com/ankur-anand/prod-todo/pkg/storage/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgreSQL provides a collection of Repository implementation over a PostgreSQL database
type PostgreSQL struct {
	// db holds connection in a pool for optimal performance
	db          *pgxpool.Pool
	userStorage postgres.UserStorage
}

// NewPostgreSQL returns an initialized PostgreSQL storage with connection pool
func NewPostgreSQL(dbURL string) (PostgreSQL, error) {
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return PostgreSQL{}, err
	}
	db, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return PostgreSQL{}, err
	}

	authPg, err := postgres.NewUserStore(db)
	if err != nil {
		return PostgreSQL{}, err
	}
	return PostgreSQL{db: db, userStorage: authPg}, nil
}

// UserStorageSQL return AUTH Repository implementation over a PostgreSQL database for User
func (p PostgreSQL) UserStorageSQL() postgres.UserStorage {
	return p.userStorage
}

// Close all the connection
func (p PostgreSQL) Close() {
	p.db.Close()
}
