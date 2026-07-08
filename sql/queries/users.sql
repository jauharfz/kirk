-- name: CreateUser :one
INSERT INTO users (
	id, created_at, updated_at,email,password,is_chirpy_red
) VALUES(
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2,
	false
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

-- name: UpdateChirpyRedStatus :one
UPDATE users
SET is_chirpy_red = $1
WHERE id = $2
RETURNING *;

