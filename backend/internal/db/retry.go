package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// WithRetryTx executes fn within a transaction, retrying up to 3 times
// if a retryable PostgreSQL error (serialization_failure, deadlock_detected) occurs.
func WithRetryTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		err = fn(tx)
		if err == nil {
			return tx.Commit()
		}
		_ = tx.Rollback()
		if !isRetryablePgError(err) || attempt == maxRetries {
			return err
		}
		time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
	}
	return nil
}

func isRetryablePgError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "40001", "40P01": // serialization_failure, deadlock_detected
			return true
		}
	}
	return false
}
