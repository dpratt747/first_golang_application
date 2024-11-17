-- +goose Up
CREATE TABLE IF NOT EXISTS user_deletes(
    user_id INT NOT NULL,
    deletion_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(user_id)
);

-- +goose Down
DROP TABLE user_deletes;