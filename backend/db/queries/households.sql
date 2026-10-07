-- name: CreateHousehold :one
INSERT INTO households (name)
VALUES (@name)
RETURNING *;

-- name: AddHouseholdMember :one
INSERT INTO household_members (household_id, user_id)
VALUES (@household_id, @user_id)
RETURNING *;

-- name: ListUserHouseholds :many
SELECT h.* FROM households h
JOIN household_members m ON m.household_id = h.id
WHERE m.user_id = @user_id
ORDER BY m.created_at;

-- name: IsHouseholdMember :one
SELECT EXISTS (
    SELECT 1 FROM household_members
    WHERE household_id = @household_id
      AND user_id = @user_id
);
