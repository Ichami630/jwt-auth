-- name: CreateUser :one
INSERT INTO "user" (name, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM "user" WHERE email = $1;

