package balance

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/model"
)

type Balance interface {
	UserBalance(ctx context.Context, userID uint64) (*model.UserBalance, error)
}
