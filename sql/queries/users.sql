-- name: CreateUser :one
INSERT INTO users (
	id, created_at, updated_at,email,password
) VALUES(
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2
)
RETURNING *;

-- name: GetUsersByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUserEmailAndPwd :one
UPDATE users
SET email = $1, password = $2
WHERE id = $3
RETURNING *;
