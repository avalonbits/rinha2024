// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package repo

import (
	"context"
	"database/sql"
)

const createTransaction = `-- name: CreateTransaction :exec
INSERT INTO Transactions (cid, tid, value, description, created_at)
       VALUES (?, ?, ?, ?, ?)
`

type CreateTransactionParams struct {
	Cid         int64
	Tid         string
	Value       int64
	Description string
	CreatedAt   string
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) error {
	_, err := q.db.ExecContext(ctx, createTransaction,
		arg.Cid,
		arg.Tid,
		arg.Value,
		arg.Description,
		arg.CreatedAt,
	)
	return err
}

const getBalance = `-- name: GetBalance :one
SELECT  SUM(value) AS balance FROM Transactions WHERE cid = ?
`

func (q *Queries) GetBalance(ctx context.Context, cid int64) (sql.NullFloat64, error) {
	row := q.db.QueryRowContext(ctx, getBalance, cid)
	var balance sql.NullFloat64
	err := row.Scan(&balance)
	return balance, err
}

const getLimit = `-- name: GetLimit :one
SELECT value FROM Limits WHERE cid = ? LIMIT 1
`

func (q *Queries) GetLimit(ctx context.Context, cid int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getLimit, cid)
	var value int64
	err := row.Scan(&value)
	return value, err
}

const transactionHistory = `-- name: TransactionHistory :many
SELECT cid, tid, value, description, created_at FROM Transactions  WHERE cid = ? ORDER BY created_at DESC LIMIT 10
`

func (q *Queries) TransactionHistory(ctx context.Context, cid int64) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, transactionHistory, cid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.Cid,
			&i.Tid,
			&i.Value,
			&i.Description,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
