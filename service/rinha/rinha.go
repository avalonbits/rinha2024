package rinha

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/avalonbits/rinha2024/storage/datastore"
	"github.com/avalonbits/rinha2024/storage/datastore/repo"
	"github.com/oklog/ulid/v2"
)

type Service struct {
	db *datastore.DB
}

func New(db *datastore.DB) *Service {
	return &Service{
		db: db,
	}
}

var (
	NotFoundErr  = fmt.Errorf("not found")
	OverLimitErr = fmt.Errorf("over limit")
)

//easyjson:json
type TransactResponse struct {
	Limit   int64 `json:"limite"`
	Balance int64 `json:"saldo"`
}

func (s *Service) Transact(
	ctx context.Context, cid, value int64, description string,
) (TransactResponse, error) {
	r := &TransactResponse{}
	now := time.Now().UTC()

	// While we don't guarantee that the tid will be used, it is better to create
	// it outside the transaction to reduce the time the transaction takes.
	tid := txID(now)

	return *r, s.db.Transaction(ctx, func(tx *datastore.DB) error {
		limit, err := tx.GetLimit(ctx, cid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return NotFoundErr
			}
			return err
		}

		bf64, err := tx.GetBalance(ctx, cid)
		if err != nil {
			return err
		}
		balance := int64(bf64.Float64) + value

		if limit+balance < 0 {
			return OverLimitErr
		}

		err = tx.CreateTransaction(ctx, repo.CreateTransactionParams{
			Cid:         cid,
			Tid:         tid,
			Value:       value,
			Description: description,
			CreatedAt:   now.Format(time.RFC3339Nano),
		})

		r.Limit = limit
		r.Balance = balance
		return nil
	})
}

//easyjson:json
type balance struct {
	Total int64  `json:"total"`
	Limit int64  `json:"limite"`
	When  string `json:"data_extracao"`
}

//easyjson:json
type transaction struct {
	Value       int64  `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
	When        string `json:"realizada_em"`
}

//easyjson:json
type AccountHistoryResponse struct {
	Balance      balance       `json:"saldo"`
	Transactions []transaction `json:"ultimas_transacoes"`
}

func (s *Service) AccountHistory(ctx context.Context, cid int64) (AccountHistoryResponse, error) {
	r := &AccountHistoryResponse{}
	now := time.Now().UTC()

	return *r, s.db.Transaction(ctx, func(tx *datastore.DB) error {
		limit, err := tx.GetLimit(ctx, cid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return NotFoundErr
			}
			return err
		}

		bf64, err := tx.GetBalance(ctx, cid)
		if err != nil {
			return err
		}
		bal := int64(bf64.Float64)

		*r = AccountHistoryResponse{
			Balance: balance{
				Limit: limit,
				Total: bal,
				When:  now.Format(time.RFC3339Nano),
			},
		}

		history, err := tx.TransactionHistory(ctx, cid)
		if err != nil {
			return err
		}
		r.Transactions = make([]transaction, 0, 0)

		for _, h := range history {
			tType := "c"
			if h.Value < 0 {
				tType = "d"
				h.Value = -h.Value
			}
			r.Transactions = append(r.Transactions, transaction{
				Value:       h.Value,
				Type:        tType,
				Description: h.Description,
				When:        h.CreatedAt,
			})
		}

		return nil

	})
}

func txID(now time.Time) string {
	entropy := ulid.Monotonic(rand.Reader, math.MaxUint64)
	return ulid.MustNew(uint64(now.UnixMilli()), entropy).String()
}
