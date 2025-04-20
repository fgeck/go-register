-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1 LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users 
WHERE id = $1 LIMIT 1;

-- name: UserExistsByEmail :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE email = $1
) AS exists;

-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = $1, email = $2, password_hash = $3
WHERE id = $4
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: DropAllUsers :exec
DELETE FROM users;
