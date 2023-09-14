-- name: Get :one
SELECT * FROM keyvalue
WHERE key = ? LIMIT 1;

-- name: Put :one
INSERT INTO keyvalue (
    key, value
) VALUES (
    ?, ?
)
RETURNING *;
