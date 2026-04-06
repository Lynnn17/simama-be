package infras

import (
	"context"

	"lms-be/shared/failure"

	"github.com/jmoiron/sqlx"
)

// TxManager handles database transactions.
// Inject this into services that need to manage transactions across multiple repositories.
type TxManager struct {
	DB *PostgresqlConn
}

// ProvideTxManager creates a new TxManager instance.
func ProvideTxManager(db *PostgresqlConn) *TxManager {
	return &TxManager{DB: db}
}

// WithTx executes the given function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func (m *TxManager) WithTx(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	e := make(chan error)
	tx, err := m.DB.Write.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	go func() {
		if err := fn(tx); err != nil {
			e <- err
			return
		}
		e <- nil
	}()

	err = <-e
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return failure.InternalError(errTx)
		}
		return err
	}
	return tx.Commit()
}
