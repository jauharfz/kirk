-- name: CreateChirp :one 
INSERT INTO chirp(
	id, created_at, updated_at, body, user_id
) VALUES (
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirp 
ORDER BY created_at ASC;

-- name: GetChirpsById :one
SELECT * FROM chirp
WHERE id = $1;

-- name: DeleteChirpById :exec
DELETE FROM chirp
WHERE user_id = $1 and id = $2;

-- name: GetChirpByUserId :many
SELECT * FROM chirp
WHERE user_id = $1
ORDER BY created_at ASC;
