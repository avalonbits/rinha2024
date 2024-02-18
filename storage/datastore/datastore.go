package datastore

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"sync"

	"github.com/avalonbits/rinha2024/storage/datastore/repo"
	"github.com/pressly/goose/v3"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*
var migrations embed.FS

const (
	readDSN  = "%s?_journal=wal&_txlock=deferred"
	writeDSN = "%s?_journal=wal&_txlock=immediate"
)

func GetDB(dbName string) (*DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf(writeDSN, dbName))
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

	wrdb, err := sql.Open("sqlite3", fmt.Sprintf(writeDSN, dbName))
	if err != nil {
		return nil, err
	}
	wrdb.SetMaxOpenConns(1)

	rddb, err := sql.Open("sqlite3", fmt.Sprintf(readDSN, dbName))
	if err != nil {
		wrdb.Close()
		return nil, err
	}
	return &DB{rddb: rddb, mu: &sync.Mutex{}, wrdb: wrdb}, nil
}

type DB struct {
	rddb *sql.DB

	mu   *sync.Mutex
	wrdb *sql.DB
}

func (db *DB) Close() error {
	return errors.Join(db.rddb.Close(), db.wrdb.Close())
}

func (db *DB) Read(ctx context.Context, f func(queries *repo.Queries) error) error {
	return db.transaction(ctx, db.rddb, f)
}

func (db *DB) Write(ctx context.Context, f func(queries *repo.Queries) error) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.transaction(ctx, db.wrdb, f)
}

func (db *DB) transaction(ctx context.Context, rdbms *sql.DB, f func(queries *repo.Queries) error) error {
	tx, err := rdbms.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}

	if err := f(repo.New(tx)); err != nil {
		rbErr := tx.Rollback()
		err = fmt.Errorf("transaction error: %w", err)

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
