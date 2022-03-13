package model

import (
	"errors"
	"github.com/shopspring/decimal"
	"github.com/theplant/luhn"
	"strconv"
	"time"
)

type OrderStatus int

const (
	OrderNew OrderStatus = iota + 1
	OrderProcessing
	OrderInvalid
	OrderProcessed
)

func (status OrderStatus) String() string {
	switch status {
	case OrderNew:
		return "NEW"
	case OrderProcessing:
		return "PROCESSING"
	case OrderProcessed:
		return "PROCESSED"
	case OrderInvalid:
		return "INVALID"
	}

	return ""
}

type Order struct {
	UserID     uint64
	Accrual    decimal.Decimal
	OrderID    string
	Status     OrderStatus
	UploadedAt time.Time
}

func (o Order) Validate() error {
	n, err := strconv.Atoi(o.OrderID)
	if err != nil {
		return errors.New("order number is not integer")
	}

	if !luhn.Valid(n) {
		return errors.New("invalid checksum for order number")
	}

	return nil
}
