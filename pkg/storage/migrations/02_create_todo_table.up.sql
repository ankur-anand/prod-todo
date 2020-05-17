CREATE TABLE IF NOT EXISTS todos (
    todo_id uuid NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL,
    title varchar(255) NOT NULL,
    content text NOT NULL,
    finished boolean NOT NULL DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);