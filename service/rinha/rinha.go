package rinha

import (
	"context"

	"github.com/avalonbits/rinha2024/storage/datastore"
)

type Service struct {
	db *datastore.DB
}

func New(db *datastore.DB) *Service {
	return &Service{
		db: db,
	}
}

//easyjson:json
type TransactResponse struct {
}

func (s *Service) Transact(
	ctx context.Context, cid, value int64, description string,
) (TransactResponse, error) {
	r := TransactResponse{}

	return r, nil
}
