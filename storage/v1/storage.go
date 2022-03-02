package v1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"github.com/soundrussian/go-practicum-diploma/storage"
)

var _ storage.Store = (*Storage)(nil)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func New() (storage.Store, error) {
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

// Log returns logger with service field set.
func (s *Storage) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.CtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceNameKey, "storage").Logger()

	return &logger
}

func runMigrations(db *sql.DB) error {
	var m *migrate.Migrate
	var driver database.Driver
	var err error

	if driver, err = postgres.WithInstance(db, &postgres.Config{}); err != nil {
		return err
	}

	if m, err = migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver); err != nil {
		return err
	}

	// golang-migrate returns ErrNoChange if there are no new migrations.
	// Ignore it.
	if err = m.Up(); !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
