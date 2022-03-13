package psql

import (
	"context"
	"database/sql"
	"github.com/shopspring/decimal"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/storage"
	"time"
)

func (s *Storage) UserBalance(ctx context.Context, userID uint64) (*model.UserBalance, error) {
	var balance model.UserBalance
	err := s.db.QueryRowContext(
		ctx,
		`SELECT COALESCE(SUM(amount), 0) AS current,
       			COALESCE(SUM(CASE WHEN amount < 0 THEN amount * -1 ELSE 0 END), 0) AS withdrawn
		   FROM transactions
		  WHERE user_id = $1
		  LIMIT 1`,
		userID,
	).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		s.Log(ctx).Err(err).Msgf("errors fetching balance for user %d", userID)
		return nil, err
	}

	return &balance, nil
}

func (s *Storage) Withdraw(ctx context.Context, userID uint64, withdrawal model.Withdrawal) (*model.Withdrawal, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.Log(ctx).Err(err).Msg("error starting transaction")
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			s.Log(ctx).Err(err).Msg("error rolling back transaction")
		}
	}()

	// Check current balance
	var currentBalance decimal.Decimal

	err = tx.QueryRowContext(
		ctx,
		`SELECT COALESCE(SUM(amount), 0)
           FROM transactions
          WHERE user_id = $1
          LIMIT 1`,
		userID,
	).Scan(&currentBalance)
	if err != nil {
		s.Log(ctx).Err(err).Msgf("failed to get current balance for user %d", userID)
		return nil, err
	}

	if currentBalance.LessThan(withdrawal.Sum) {
		s.Log(ctx).Info().Msgf("user's %d current balance of %d is less than withdrawal amount of %d", userID, currentBalance, withdrawal.Sum)
		return nil, storage.ErrNotEnoughBalance
	}

	now := time.Now()
	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO transactions(user_id, order_id, amount, created_at)
              VALUES ($1, $2, $3, $4)`,
		userID, withdrawal.Order, withdrawal.Sum.Mul(decimal.NewFromInt(-1)), now,
	)
	if err != nil {
		s.Log(ctx).Err(err).Msg("error saving withdrawal to DB")
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		s.Log(ctx).Err(err).Msg("error committing transaction")
		return nil, err
	}

	return &model.Withdrawal{
		Order:       withdrawal.Order,
		Sum:         withdrawal.Sum,
		ProcessedAt: now,
	}, nil
}

func (s *Storage) UserWithdrawals(ctx context.Context, userID uint64) ([]model.Withdrawal, error) {
	result := make([]model.Withdrawal, 0)
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT order_id,
                amount * -1,
                created_at
		   FROM transactions
		  WHERE user_id = $1
            AND amount < 0
				`,
		userID)
	if err != nil {
		s.Log(ctx).Err(err).Msgf("failed to fetch withdrawals for user_id %d", userID)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		record := model.Withdrawal{}
		if err := rows.Scan(&record.Order, &record.Sum, &record.ProcessedAt); err != nil {
			s.Log(ctx).Err(err).Msg("failed to scan row")
			return nil, err
		}
		result = append(result, record)
	}

	if err := rows.Err(); err != nil {
		s.Log(ctx).Err(err).Msg("error reading rows")
		return nil, err
	}

	return result, nil
}
