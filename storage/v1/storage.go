package v1

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/soundrussian/go-practicum-diploma/model"
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

func (s *Storage) CreateUser(login string, password string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
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

	return m.Up()
}
