-- name: CreateToken :one

INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;

-- name: RetrieveToken :one

SELECT
    token,
    expires_at,
    revoked_at,
    user_id
FROM refresh_tokens
WHERE token = $1;

-- name: GetUserFromRefreshToken :one

SELECT 
    user_id
FROM refresh_tokens
WHERE token = $1;

-- name: RevokeToken :exec

UPDATE refresh_tokens
SET 
    revoked_at = $2,
    updated_at = $2
WHERE user_id = $1;