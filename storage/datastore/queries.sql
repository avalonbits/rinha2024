-- name: CreateTransaction :exec
INSERT INTO Transactions (cid, tid, value, description, created_at)
       VALUES (?, ?, ?, ?, ?);

-- name: GetBalance :one
SELECT  SUM(value) AS balance FROM Transactions WHERE cid = ?;

-- name: GetLimit :one
SELECT value FROM Limits WHERE cid = ? LIMIT 1;

-- name: TransactionHistory :many
SELECT * FROM Transactions  WHERE cid = ? ORDER BY created_at DESC LIMIT 10;
