CREATE TABLE IF NOT EXISTS todos (
    todo_id uuid NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    user_id uuid NOT NULL,
    title varchar(255) NOT NULL,
    content text NOT NULL,
    finished boolean NOT NULL DEFAULT FALSE,
    CONSTRAINT todo_pk PRIMARY KEY(todo_id, created_at),
    CONSTRAINT todo_fk FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);