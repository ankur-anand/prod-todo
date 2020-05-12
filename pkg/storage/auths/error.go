package auths

import "errors"

var (
	// ErrUserNotFound indicates no user associated with either ID or emailID
	ErrUserNotFound = errors.New("no user found")
	// ErrInsertCommand indicates error with insert query operation
	ErrInsertCommand = errors.New("insert command operation")
)
