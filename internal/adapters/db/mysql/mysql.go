package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"testcode/test3/pkg/logger"
	"testcode/test3/pkg/mysql"
)

const _defaultEntityCap = 64

type baseStorage struct {
	db *mysql.Mysql
}

type transactor struct {
	baseStorage
	log *logger.Logger
}

func NewTransactor(log *logger.Logger, db *mysql.Mysql) *transactor {
	return &transactor{
		baseStorage{db},
		log,
	}
}

type txKey struct{}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func (r *baseStorage) Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	tx := extractTx(ctx)
	if tx != nil {
		return (*tx).ExecContext(ctx, sql, args...)
	}
	return r.db.Pool.ExecContext(ctx, sql, args...)
}

func (r *baseStorage) Query(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error) {
	tx := extractTx(ctx)
	if tx != nil {
		return (*tx).QueryContext(ctx, sql, args...)
	}
	return r.db.Pool.QueryContext(ctx, sql, args...)
}

// WithinTransaction runs function within transaction
//
// The transaction commits when function were finished without error
func (r *transactor) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) (err error) {
	// begin transaction
	tx, err := r.db.Pool.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		// finalize transaction on panic, etc.
		errTx := tx.Rollback()
		if errTx != nil && !errors.Is(errTx, sql.ErrTxDone) {
			r.log.Info("close transaction: %v", errTx)
			err = errTx
		}
	}()

	// run callback
	err = tFunc(injectTx(ctx, tx))
	if err != nil {
		// if error, rollback
		_ = tx.Rollback()
		r.log.Info("rollback transaction: %v", err)
		return err
	}

	// if no error, commit
	return tx.Commit()
}
