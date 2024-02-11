-- name: CreateTransaction :exec
INSERT INTO Transactions (cid, tid, value, description)
       VALUES (?, ?, ?, ?);

-- name: GetBalance :one
SELECT  L.Value value,
        (SELECT SUM(value) FROM Transactions T WHERE T.cid = L.cid) balance
FROM Limits L WHERE L.cid = ?;

-- name: TransactionHistory :many
SELECT tid, value, description FROM Transactions  WHERE cid = ? ORDER BY tid DESC LIMIT 10;
