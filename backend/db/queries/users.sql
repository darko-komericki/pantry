-- name: CreateUser :one
INSERT INTO users (email, password_hash, display_name)
VALUES (@email, @password_hash, @display_name)
RETURNING *;

-- name: GetUserByEmail :one
-- lower() on both sides matches the users_email_lower_key index.
SELECT * FROM users
WHERE lower(email) = lower(@email);

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = @id;
