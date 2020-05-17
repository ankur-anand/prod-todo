package todo

import (
	"github.com/ankur-anand/prod-todo/pkg/storage/serr"
	"github.com/jackc/pgx/v4"

	"github.com/ankur-anand/prod-todo/pkg/domain"
)

// todoIterator is a domain.TodoIterator implementation for the postgres.
type todoIterator struct {
	rows        pgx.Rows
	lastErr     error
	latchedTodo domain.TodoModel
	query       string
}

// Next implements domain.TodoIterator
func (i *todoIterator) Next() bool {
	if i.lastErr != nil || !i.rows.Next() {
		return false
	}

	var todo domain.TodoModel
	i.lastErr = i.rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Content, &todo.Finished)
	if i.lastErr != nil {
		return false
	}

	i.latchedTodo = todo
	return true
}

// Error implements graph.LinkIterator.
func (i *todoIterator) Error() error {
	switch i.lastErr {
	case nil:
		return nil
	case pgx.ErrNoRows:
		return serr.NewQueryError(i.query, ErrNoMoreRow, i.lastErr.Error())
	default:
		// something else went wrong,
		return serr.NewQueryError(i.query, i.lastErr, i.lastErr.Error())
	}
}

// Close implements domain.TodoIterator.
func (i *todoIterator) Close() error {
	i.rows.Close()
	return nil
}

// Todo implements domain.TodoIterator.
func (i *todoIterator) Todo() domain.TodoModel {
	return i.latchedTodo
}
