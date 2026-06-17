-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdatePassword :one
UPDATE users
SET password_hash = $1
WHERE id = $2
RETURNING *;

-- name: SetPasswordResetToken :one
UPDATE users
SET password_reset_token      = $1,
    password_reset_expires_at = $2
WHERE email = $3
RETURNING *;

-- name: GetUserByPasswordResetToken :one
SELECT * FROM users
WHERE password_reset_token = $1
  AND password_reset_expires_at > NOW();

-- name: ResetPassword :one
UPDATE users
SET password_hash             = $1,
    password_reset_token      = NULL,
    password_reset_expires_at = NULL
WHERE id = $2
RETURNING *;

-- name: UpdateTier :one
UPDATE users
SET tier = $1
WHERE id = $2
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
