package infras

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type contextKey string

const transactionKey contextKey = "transaction"

// ContextWithTx creates a new context with the given transaction
func ContextWithTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, transactionKey, tx)
}

// TxFromContext returns the transaction from the context, or nil if not found
func TxFromContext(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(transactionKey).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}
