package psql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/storage"
)

func (s *Storage) CreateUser(ctx context.Context, login string, password string) (*model.User, error) {
	var recordID uint64
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO users(login, encrypted_password)
				VALUES ($1, $2)
		 RETURNING id`,
		login, password,
	).Scan(&recordID)
	if err != nil {
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

	err := s.db.QueryRowContext(
		ctx,
		`SELECT id, login, encrypted_password
           FROM users
          WHERE login = $1
          LIMIT 1`,
		login,
	).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		s.Log(ctx).Err(err).Msgf("error fetching user %s", login)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
