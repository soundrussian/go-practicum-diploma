package balance

import "errors"

var (
	ErrNotEnoughBalance = errors.New("not enough balance")
	ErrInvalidSum       = errors.New("withdrawal sum must be greater than zero")
	ErrInvalidOrder     = errors.New("invalid order number")
)
