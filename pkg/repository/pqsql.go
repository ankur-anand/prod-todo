package repository

import (
	"context"

	"github.com/ankur-anand/prod-app/pkg/repository/authstorage"

	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgreSQL provides a collection of Repository implementation over a PostgreSQL database
type PostgreSQL struct {
	// db holds connection in a pool for optimal performance
	db      *pgxpool.Pool
	authSQL authstorage.PgSQL
}

// NewPostgreSQL returns an initialized PostgreSQL repository with connection pool
func NewPostgreSQL(dbURL string) (PostgreSQL, error) {
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return PostgreSQL{}, err
	}
	db, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return PostgreSQL{}, err
	}

	authPg, err := authstorage.NewAuthPgSQL(db)
	if err != nil {
		return PostgreSQL{}, err
	}
	return PostgreSQL{db: db, authSQL: authPg}, nil
}

// AuthSQL return AUTH Repository implementation over a PostgreSQL database
func (p PostgreSQL) AuthSQL() authstorage.PgSQL {
	return p.authSQL
}
