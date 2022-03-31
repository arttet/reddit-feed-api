package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type WithTxFunc func(ctx context.Context, tx *sqlx.Tx) error

func WithTx(ctx context.Context, db *sqlx.DB, fn WithTxFunc) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "db.BeginTxx()")
	}

	if err = fn(ctx, tx); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errors.Wrap(err, "Tx.Rollback")
		}

		return errors.Wrap(err, "Tx.WithTxFunc")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "Tx.Commit")
	}

	return nil
}
