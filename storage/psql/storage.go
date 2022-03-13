package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/soundrussian/go-practicum-diploma/storage"
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

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://storage/psql/db/migrations", "postgres", driver)
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
