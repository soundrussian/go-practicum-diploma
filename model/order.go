package model

import (
	"errors"
	"github.com/theplant/luhn"
	"strconv"
)

var (
	ErrNotNum          = errors.New("order number is not integer")
	ErrInvalidChecksum = errors.New("invalid checksum for order number")
)

type Order struct {
	Num string
}

func (o Order) Validate() error {
	var n int
	var err error

	if n, err = strconv.Atoi(o.Num); err != nil {
		return ErrNotNum
	}

	if !luhn.Valid(n) {
		return ErrInvalidChecksum
	}

	return nil
}
