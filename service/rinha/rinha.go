package rinha

import "github.com/avalonbits/rinha2024/storage/datastore"

type Service struct {
	db *datastore.DB
}

func New(db *datastore.DB) *Service {
	return &Service{
		db: db,
	}
}
