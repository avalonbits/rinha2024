-- name: CreateTransaction :exec
INSERT INTO Transaction (cid, tid, value, description, created_at)
       VALUES (?, ?, ?, ?, ?);

-- name: GetBalance :one
SELECT  FLOOR(SUM(value)) AS balance FROM Transaction WHERE cid = ?;

-- name: GetLimit :one
SELECT value FROM Limit WHERE cid = ? LIMIT 1;

-- name: TransactionHistory :many
SELECT * FROM Transaction  WHERE cid = ? ORDER BY created_at DESC;
