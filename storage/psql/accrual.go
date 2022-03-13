package psql

import (
	"context"
	"database/sql"
	"github.com/shopspring/decimal"
	"github.com/soundrussian/go-practicum-diploma/model"
	"time"
)

func (s *Storage) AddAccrual(ctx context.Context, orderID string, status model.OrderStatus, accrual decimal.Decimal) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.Log(ctx).Err(err).Msg("error starting transaction")
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			s.Log(ctx).Err(err).Msg("error rolling back transaction")
		}
	}()

	var userID uint64

	// Only orders with status processing can be awarded accrual
	err = tx.QueryRowContext(
		ctx,
		`UPDATE orders
            SET accrual = $1,
                status = $2::integer
          WHERE order_id = $3
            AND status = $4
         RETURNING user_id`,
		accrual, status, orderID, model.OrderProcessing,
	).Scan(&userID)
	if err != nil {
		s.Log(ctx).Err(err).Msgf("failed to update accrual to %d for order %s", accrual, orderID)
		return err
	}

	now := time.Now()

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO transactions (user_id, order_id, amount, created_at)
              VALUES ($1, $2, $3, $4)`,
		userID, orderID, accrual, now,
	)
	if err != nil {
		s.Log(ctx).Err(err).Msgf("failed to save transaction for user %d", userID)
		return err
	}
	s.Log(ctx).Debug().Msgf("inserted %d, %s, %s", userID, orderID, accrual)

	if err := tx.Commit(); err != nil && err != sql.ErrTxDone {
		s.Log(ctx).Err(err).Msg("error committing transaction")
		return err
	}

	return nil
}
