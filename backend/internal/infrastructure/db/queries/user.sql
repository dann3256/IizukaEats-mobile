-- name: CreateUser :one
INSERT INTO users (name,email, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT id,name, email, password_hash
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT name, email, password_hash
FROM users
WHERE id = $1;




-- name: UpdateUser :exec
UPDATE users
SET name = $2, email=$3, updated_at = NOW()
WHERE id = $1;


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;