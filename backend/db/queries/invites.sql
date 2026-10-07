-- name: CreateInvite :one
INSERT INTO household_invites (household_id, code, created_by, expires_at)
VALUES (@household_id, @code, @created_by, @expires_at)
RETURNING *;

-- name: ConsumeInvite :one
-- Checks and uses the invite in one statement. If two sign-ups send the same
-- code at once, only one gets the row back; the other gets no rows.
DELETE FROM household_invites
WHERE code = @code
  AND expires_at > now()
RETURNING *;

-- name: DeleteExpiredInvites :execrows
DELETE FROM household_invites
WHERE expires_at <= now();
