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

const (
	filterPostgresTrue  = "TRUE"
	filterPostgresFalse = "FALSE"
)

// TodoStorage provides a ToDO Storage implementation over a PostgreSQL database
type TodoStorage struct {
	// db holds connection in a pool for optimal performance
	db *pgxpool.Pool
}

func (t TodoStorage) FindOneTodo(ctx context.Context, id uuid.UUID) (pkg.TodoModel, error) {
	var todo pkg.TodoModel
	err := t.db.QueryRow(ctx, findTodoByIDQuery, id).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Content, &todo.Finished)
	switch err {
	case nil:
		return todo, nil
	case pgx.ErrNoRows:
		return pkg.NilTodoModel, serror.NewQueryError(findTodoByIDQuery, serror.ErrTodoNotFound, err.Error())
	default:
		return pkg.NilTodoModel, serror.NewQueryError(findTodoByIDQuery, err, err.Error())
	}
}

func (t TodoStorage) FindAllTodoOfUser(ctx context.Context, userID uuid.UUID, filter pkg.TodoFilter) ([]pkg.TodoModel, error) {
	query, err := getFilterValue(filter)
	if err != nil {
		return nil, err
	}
	rows, err := t.db.Query(ctx, query, userID)
	if err != nil {
		return nil, serror.NewQueryError(query, err, err.Error())
	}
}

func (t TodoStorage) UpdateOne(ctx context.Context, todo pkg.TodoModel) error {
	panic("implement me")
}

func (t TodoStorage) InsertOne(ctx context.Context, todo pkg.TodoModel) (uuid.UUID, error) {
	panic("implement me")
}

func (t TodoStorage) DeleteOne(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}

func getFilterValue(filter pkg.TodoFilter) (string, error) {
	switch filter {
	case pkg.NilFilter:
		return findAllTodoByUser, nil
	case pkg.Finished:
		return fmt.Sprintf(findAllTodoByUserWithFinishedFilter, "TRUE"), nil
	case pkg.UnFinished:
		return fmt.Sprintf(findAllTodoByUserWithFinishedFilter, "FALSE"), nil
	default:
		return "", fmt.Errorf("unsupported filter")
	}
}
