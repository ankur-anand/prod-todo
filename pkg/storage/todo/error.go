package todo

import "errors"

var (
	// ErrTodoNotFound indicates no todo associated with either todoID or userID
	ErrTodoNotFound = errors.New("no todo found")
	// ErrNoMoreRow indicates there is no rows left to read
	ErrNoMoreRow = errors.New("no more rows available")
	// ErrInsertCommand indicates error with insert query operation
	ErrInsertCommand = errors.New("insert command operation")
	// ErrDeleteCommand indicates error with insert query operation
	ErrDeleteCommand = errors.New("delete command operation")
)
