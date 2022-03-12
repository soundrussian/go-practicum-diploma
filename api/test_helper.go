package api

import (
	"context"
	"github.com/soundrussian/go-practicum-diploma/model"
	"github.com/soundrussian/go-practicum-diploma/service/auth/v1"
)

func token(userID uint64) string {
	a := &v1.Auth{}
	t, _ := a.AuthToken(context.Background(), &model.User{ID: userID})
	return *t
}
