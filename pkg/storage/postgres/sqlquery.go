package postgres

var (
	// SQL Query
	findUserByIDQuery    = "SELECT user_id, email_id, password_hash, first_name, last_name, user_name FROM users WHERE user_id=$1"
	findUserByEmailQuery = "SELECT user_id, email_id, password_hash, first_name, last_name, user_name FROM users WHERE email_id=$1"
	updateUserQuery      = "UPDATE users SET email_id = $2, password_hash = $3, first_name = $4, last_name = $5, user_name = $6 WHERE user_id = $1"
	storeUserQuery       = `
INSERT INTO users (user_id, email_id, password_hash, first_name, last_name, user_name) VALUES ($1, $2, $3, $4, $5, $6)
`
)

var (

	// SQL Query
	findAllTodo                         = "SELECT todo_id, user_id, title, content, finished FROM todos ORDER BY created_at DESC LIMIT 50"
	findTodoByIDQuery                   = "SELECT * FROM todos WHERE todo_id=$1"
	findAllTodoByUser                   = "SELECT * FROM todos WHERE user_id=$1 ORDER BY created_at DESC LIMIT 50"
	findAllTodoByUserWithFinishedFilter = "SELECT * FROM todos WHERE user_id=$1 AND finished = %s ORDER BY created_at DESC LIMIT 50"

	storeTodoQuery = `
INSERT INTO todos (user_id, email_id, password_hash, first_name, last_name, user_name) VALUES ($1, $2, $3, $4, $5, $6)
`
	deleteTodoByID = "DROP FROM todos WHERE todo_id=$1"
)
