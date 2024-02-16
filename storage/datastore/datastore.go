package datastore

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/avalonbits/rinha2024/storage/datastore/repo"
	"github.com/pressly/goose/v3"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*
var migrations embed.FS

func GetWriteDB(dbURL string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		return nil, err
	}

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		db.Close()
		return nil, err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		db.Close()
		return nil, err
	}
	db.Close()

	db, err = sql.Open("sqlite3", dbURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	return &DB{Queries: repo.New(db), rdbms: db}, nil
}

func GetReadDB(dbURL string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	return &DB{Queries: repo.New(db), rdbms: db}, nil
}

type DB struct {
	*repo.Queries
	rdbms *sql.DB
}

func (db *DB) Close() error {
	return db.rdbms.Close()
}

func (db *DB) RDBMS() *sql.DB {
	return db.rdbms
}

func (db *DB) Transaction(ctx context.Context, f func(*DB) error) error {
	tx, err := db.rdbms.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}

	txdb := &DB{
		rdbms:   db.rdbms,
		Queries: db.Queries.WithTx(tx),
	}
	if err := f(txdb); err != nil {
		err = fmt.Errorf("transaction error: %w", err)
		rbErr := tx.Rollback()

		if rbErr != nil {
			err = errors.Join(err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

func NoRows(err error) bool {
	return err != nil && errors.Is(err, sql.ErrNoRows)
}
