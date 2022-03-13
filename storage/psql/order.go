package psql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/storage"
	"time"
)

func (s *Storage) AcceptOrder(ctx context.Context, userID uint64, orderID string) (*model.Order, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.Log(ctx).Err(err).Msg("error starting transaction")
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			s.Log(ctx).Err(err).Msg("error rolling back transaction")
		}
	}()

	if exists := s.sameOrAnotherUser(ctx, tx, orderID, userID); exists != nil {
		return nil, exists
	}

	var order model.Order
	now := time.Now()

	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO orders (order_id, user_id, uploaded_at) 
			  VALUES ($1, $2, $3)
		 RETURNING order_id, user_id, status, uploaded_at`,
		orderID, userID, now,
	).Scan(&order.OrderID, &order.UserID, &order.Status, &order.UploadedAt)
	if err != nil {
		s.Log(ctx).Err(err).Msgf("error while saving order <%s> for user %d", orderID, userID)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		s.Log(ctx).Err(err).Msg("error while committing transaction")
		return nil, err
	}

	return &order, nil
}

func (s *Storage) OrdersWithStatus(ctx context.Context, status model.OrderStatus, limit int) ([]string, error) {
	result := make([]string, 0)

	rows, err := s.db.QueryContext(
		ctx,
		`SELECT order_id
		   FROM orders
		  WHERE status = $1::integer
		  LIMIT $2`,
		int(status), limit,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Log(ctx).Info().Msgf("no orders with status %d", status)
			return result, nil
		}
		s.Log(ctx).Err(err).Msgf("failed to fetch orders with status %d", status)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var orderID string
		if err := rows.Scan(&orderID); err != nil {
			s.Log(ctx).Err(err).Msg("failed to scan row")
			return nil, err
		}
		result = append(result, orderID)
	}

	if err := rows.Err(); err != nil {
		s.Log(ctx).Err(err).Msg("error reading rows")
		return nil, err
	}

	return result, nil
}

func (s *Storage) UpdateOrderStatus(ctx context.Context, orderID string, status model.OrderStatus) error {
	_, err := s.db.ExecContext(
		ctx,
		`UPDATE orders
		 SET status = $1::integer
		 WHERE order_id = $2`,
		int(status), orderID,
	)
	if err != nil {
		s.Log(ctx).Err(err).Msgf("failed updating order <%s> with status %d", orderID, status)
		return err
	}

	return nil
}

func (s *Storage) UserOrders(ctx context.Context, userID uint64) ([]model.Order, error) {
	result := make([]model.Order, 0)

	rows, err := s.db.QueryContext(
		ctx,
		`SELECT order_id, user_id, accrual, status, uploaded_at
		   FROM orders
		  WHERE user_id = $1
          ORDER BY uploaded_at DESC`,
		userID)
	if err != nil {
		s.Log(ctx).Err(err).Msgf("failed to fetch orders for user_id %d", userID)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		record := model.Order{}
		if err := rows.Scan(&record.OrderID, &record.UserID, &record.Accrual, &record.Status, &record.UploadedAt); err != nil {
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

func (s *Storage) sameOrAnotherUser(ctx context.Context, tx *sql.Tx, orderID string, currentUserID uint64) error {
	var existingUser uint64
	err := tx.QueryRowContext(
		ctx,
		`SELECT user_id
           FROM orders
          WHERE order_id = $1
         LIMIT 1`,
		orderID,
	).Scan(&existingUser)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.Log(ctx).Err(err).Msgf("error determining which user uploaded order <%s>", orderID)
		return err
	}

	if existingUser == 0 {
		return nil
	}

	if existingUser == currentUserID {
		return storage.ErrOrderExistsSameUser
	}

	return storage.ErrOrderExistsAnotherUser
}
