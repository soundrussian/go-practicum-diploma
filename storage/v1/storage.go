package v1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"github.com/soundrussian/go-practicum-diploma/storage"
	"math"
	"time"
)

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	db *sql.DB
}

func New() (storage.Storage, error) {
	if databaseConnection == nil {
		return nil, errors.New("databaseConnection config is not set")
	}

	db, err := sql.Open("pgx", *databaseConnection)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database connection: %w", err)
	}

	if err = runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	store := Storage{db: db}

	return &store, nil
}

func (s *Storage) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *Storage) AddAccrual(ctx context.Context, orderID string, status model.OrderStatus, accrual float64) error {
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

	// Convert accrual to int
	accrualSum := int(math.Round(accrual * 100))
	var userID uint64

	if err = tx.QueryRowContext(ctx,
		`UPDATE orders SET accrual = $1, status = $2::integer WHERE order_id = $3 RETURNING user_id`,
		accrualSum, status, orderID).Scan(&userID); err != nil {
		s.Log(ctx).Err(err).Msgf("failed to update accrual to %d for order %s", accrualSum, orderID)
		return err
	}

	now := time.Now()

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO transactions (user_id, order_id, amount, created_at) VALUES ($1, $2, $3, $4)`,
		userID, orderID, accrualSum, now); err != nil {
		s.Log(ctx).Err(err).Msgf("failed to save transaction for user %d", userID)
		return err
	}
	s.Log(ctx).Debug().Msgf("inserted %d, %s, %d", userID, orderID, accrualSum)

	if err := tx.Commit(); err != nil && err != sql.ErrTxDone {
		s.Log(ctx).Err(err).Msg("error committing transaction")
		return err
	}

	return nil
}

func (s *Storage) UpdateOrderStatus(ctx context.Context, orderID string, status model.OrderStatus) error {
	if _, err := s.db.ExecContext(ctx,
		`UPDATE orders
		SET status = $1::integer
		WHERE order_id = $2
		`, int(status), orderID); err != nil {
		s.Log(ctx).Err(err).Msgf("failed updating order <%s> with status %d", orderID, status)
		return err
	}

	return nil
}

func (s *Storage) OrdersWithStatus(ctx context.Context, status model.OrderStatus, limit int) ([]string, error) {
	result := make([]string, 0)

	rows, err := s.db.QueryContext(ctx,
		`SELECT order_id
				FROM orders
				WHERE status = $1::integer
				LIMIT $2
				`,
		int(status), limit)
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

func (s *Storage) CreateUser(ctx context.Context, login string, password string) (*model.User, error) {
	var recordID uint64
	if err := s.db.QueryRowContext(
		ctx,
		"INSERT INTO users(login, encrypted_password) VALUES ($1, $2) RETURNING id",
		login, password,
	).Scan(&recordID); err != nil {
		var pgError pgx.PgError
		if errors.As(err, &pgError) && pgerrcode.UniqueViolation == pgError.Code {
			s.Log(ctx).Err(err).Msgf("user with login %s already exists", login)
			return nil, storage.ErrLoginAlreadyExists
		}
		s.Log(ctx).Err(err).Msg("failed to create user")
		return nil, err
	}

	user := &model.User{
		ID:    recordID,
		Login: login,
	}

	return user, nil
}

func (s *Storage) FetchUser(ctx context.Context, login string) (*model.User, error) {
	var user model.User

	if err := s.db.QueryRowContext(ctx,
		`SELECT id, login, encrypted_password FROM users WHERE login = $1 LIMIT 1`,
		login,
	).Scan(&user.ID, &user.Login, &user.Password); err != nil {
		s.Log(ctx).Err(err).Msgf("error fetching user %s", login)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *Storage) UserBalance(ctx context.Context, userID uint64) (*model.UserBalance, error) {
	var balance model.UserBalance
	if err := s.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(amount), 0) AS current,
       				  COALESCE(SUM(CASE WHEN amount < 0 THEN amount * -1 ELSE 0 END), 0) AS withdrawn
				FROM transactions
				WHERE user_id = $1
				LIMIT 1`,
		userID,
	).Scan(&balance.Current, &balance.Withdrawn); err != nil {
		s.Log(ctx).Err(err).Msgf("errors fetching balance for user %d", userID)
		return nil, err
	}

	// Convert kopecks to rubles
	balance.Current = balance.Current / 100.0
	balance.Withdrawn = balance.Withdrawn / 100.0

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
	var currentBalance int
	// Convert withdrawal to int
	withdrawSum := int(math.Round(withdrawal.Sum * 100))

	if err := tx.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE user_id = $1 LIMIT 1`,
		userID).Scan(&currentBalance); err != nil {
		s.Log(ctx).Err(err).Msgf("failed to get current balance for user %d", userID)
		return nil, err
	}

	if currentBalance < withdrawSum {
		s.Log(ctx).Info().Msgf("user's %d current balance of %d is less than withdrawal amount of %d", userID, currentBalance, withdrawSum)
		return nil, storage.ErrNotEnoughBalance
	}

	now := time.Now()
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO transactions(user_id, order_id, amount, created_at) VALUES ($1, $2, $3, $4)`,
		userID, withdrawal.Order, withdrawSum*-1, now,
	); err != nil {
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
	rows, err := s.db.QueryContext(ctx,
		`SELECT order_id AS "order", amount * -1 AS sum, created_at AS processed_at
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
		// Convert sum to float
		record.Sum = record.Sum / 100
		result = append(result, record)
	}

	if err := rows.Err(); err != nil {
		s.Log(ctx).Err(err).Msg("error reading rows")
		return nil, err
	}

	return result, nil
}

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

	if err := tx.QueryRowContext(ctx,
		`INSERT INTO orders (order_id, user_id, uploaded_at) 
				 VALUES ($1, $2, $3)
				 RETURNING order_id, user_id, status, uploaded_at`,
		orderID, userID, now,
	).Scan(&order.OrderID, &order.UserID, &order.Status, &order.UploadedAt); err != nil {
		s.Log(ctx).Err(err).Msgf("error while saving order <%s> for user %d", orderID, userID)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		s.Log(ctx).Err(err).Msg("error while committing transaction")
		return nil, err
	}

	return &order, nil
}

func (s *Storage) UserOrders(ctx context.Context, userID uint64) ([]model.Order, error) {
	result := make([]model.Order, 0)
	rows, err := s.db.QueryContext(ctx,
		`SELECT order_id, user_id, accrual, status, uploaded_at
				FROM orders
				WHERE user_id = $1
                ORDER BY uploaded_at DESC
				`,
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
		// Convert kopecks to rubles
		record.Accrual = record.Accrual / 100.0

		result = append(result, record)
	}

	if err := rows.Err(); err != nil {
		s.Log(ctx).Err(err).Msg("error reading rows")
		return nil, err
	}

	return result, nil
}

// Log returns logger with service field set.
func (s *Storage) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.CtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceNameKey, "storage").Logger()

	return &logger
}

func (s *Storage) sameOrAnotherUser(ctx context.Context, tx *sql.Tx, orderID string, currentUserID uint64) error {
	var existingUser uint64
	if err := tx.QueryRowContext(ctx,
		`SELECT user_id FROM orders WHERE order_id = $1 LIMIT 1`,
		orderID).Scan(&existingUser); err != nil && !errors.Is(err, sql.ErrNoRows) {
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

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	// golang-migrate returns ErrNoChange if there are no new migrations.
	// Ignore it.
	if err := m.Up(); !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
