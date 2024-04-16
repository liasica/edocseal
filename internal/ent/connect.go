// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-16, by liasica

package ent

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/liasica/edocseal/internal/ent/migrate"
)

var db *Client

func NewDatabase() *Client {
	return db
}

func CreateDatabase(dsn string, debug bool) (err error) {
	db, err = Open("postgres", dsn)
	if err != nil {
		return
	}
	db.debug = debug

	ctx := context.Background()

	if err = db.Schema.Create(ctx); err != nil {
		return err
	}

	return db.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
		migrate.WithForeignKeys(false),
	)
}

type TxFunc func(tx *Tx) (err error)

func WithTx(ctx context.Context, fn TxFunc) error {
	tx, err := db.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()
	if err = fn(tx); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			err = fmt.Errorf("rolling back transaction: %w", txErr)
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
