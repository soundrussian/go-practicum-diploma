package api

import (
	"context"
	v1 "github.com/soundrussian/go-practicum-diploma/auth/v1"
	"github.com/soundrussian/go-practicum-diploma/model"
)

func token(userID uint64) string {
	a := &v1.Auth{}
	t, _ := a.AuthToken(context.Background(), &model.User{ID: userID})
	return *t
}
