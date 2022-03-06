package model

import (
	"errors"
	"github.com/theplant/luhn"
	"strconv"
	"time"
)

var (
	ErrNotNum          = errors.New("order number is not integer")
	ErrInvalidChecksum = errors.New("invalid checksum for order number")
)

type OrderStatus int

const (
	New OrderStatus = iota
	Processing
	Invalid
	Processed
)

type Order struct {
	UserID     uint64
	Accrual    int
	OrderID    string
	Status     OrderStatus
	UploadedAt time.Time
}

func (o Order) Validate() error {
	var n int
	var err error

	if n, err = strconv.Atoi(o.OrderID); err != nil {
		return ErrNotNum
	}

	if !luhn.Valid(n) {
		return ErrInvalidChecksum
	}

	return nil
}
