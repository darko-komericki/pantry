-- name: CreateSession :one
INSERT INTO sessions (user_id, token_hash, expires_at)
VALUES (@user_id, @token_hash, @expires_at)
RETURNING *;

-- name: GetActiveSessionByTokenHash :one
SELECT * FROM sessions
WHERE token_hash = @token_hash
  AND expires_at > now();

-- name: DeleteSessionByTokenHash :exec
DELETE FROM sessions
WHERE token_hash = @token_hash;

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = @user_id;

-- name: DeleteExpiredSessions :execrows
DELETE FROM sessions
WHERE expires_at <= now();
