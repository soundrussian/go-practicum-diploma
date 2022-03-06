package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/soundrussian/go-practicum-diploma/auth"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/pkg/logging"
	"github.com/soundrussian/go-practicum-diploma/storage"
	"golang.org/x/crypto/bcrypt"
)

var _ auth.Auth = (*Auth)(nil)

type Auth struct {
	storage storage.Storage
}

func New(storage storage.Storage) (*Auth, error) {
	if secretKey == nil {
		return nil, errors.New("secretKey has not been initialized")
	}

	if TokenAuth == nil {
		return nil, errors.New("TokenAuth has not been initialized")
	}

	if storage == nil {
		return nil, errors.New("nil storage passed to Auth service constructor")
	}

	auth := &Auth{storage: storage}

	return auth, nil
}

func (a *Auth) Register(ctx context.Context, login string, password string) (*model.User, error) {
	var hashedPassword string
	var user *model.User
	var err error

	if hashedPassword, err = a.hashedPassword(ctx, password); err != nil {
		return nil, auth.ErrRegistrationInternalError
	}

	if user, err = a.storage.CreateUser(ctx, login, hashedPassword); err != nil {
		a.Log(ctx).Err(err).Msg("failed to create user")
		return nil, fmt.Errorf("failed to store user: %w", err)
	}

	return user, nil
}

func (a *Auth) Authenticate(ctx context.Context, login string, password string) (*model.User, error) {
	var user *model.User
	var err error

	if user, err = a.storage.FetchUser(ctx, login); err != nil {
		a.Log(ctx).Err(err).Msg("errors fetching user for authentication")
		if errors.Is(err, storage.ErrNotFound) {
			return nil, auth.ErrUserNotFound
		}
		return nil, auth.ErrAuthenticateInternalError
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+*secretKey)); err != nil {
		a.Log(ctx).Err(err).Msgf("provided password %s is not valid for user %s", password, login)
		return nil, auth.ErrPasswordIncorrect
	}

	return user, nil
}

func (a *Auth) AuthToken(ctx context.Context, user *model.User) (*string, error) {
	var tokenString string
	var err error

	if _, tokenString, err = TokenAuth.Encode(map[string]interface{}{"user_id": user.ID}); err != nil {
		return nil, err
	}

	return &tokenString, nil
}

// Log returns logger with service field set.
func (a Auth) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.CtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceNameKey, "auth").Logger()

	return &logger
}

func (a Auth) hashedPassword(ctx context.Context, password string) (string, error) {
	var result []byte
	var err error

	if result, err = bcrypt.GenerateFromPassword([]byte(password+*secretKey), bcrypt.DefaultCost); err != nil {
		a.Log(ctx).Err(err).Msg("failed to generate hashed password")
		return "", err
	}

	return string(result), nil
}
