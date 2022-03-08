package model

import (
	"encoding/json"
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
	OrderNew OrderStatus = iota
	OrderProcessing
	OrderInvalid
	OrderProcessed
)

type Order struct {
	UserID     uint64      `json:"-"`
	Accrual    float64     `json:"accrual,omitempty"`
	OrderID    string      `json:"number"`
	Status     OrderStatus `json:"-"`
	UploadedAt time.Time   `json:"-"`
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

// MarshalJSON is used to convert ProcessedAt to required RFC3339 format on marshal
func (o *Order) MarshalJSON() ([]byte, error) {
	var status string

	switch o.Status {
	case OrderNew:
		status = "NEW"
	case OrderProcessing:
		status = "PROCESSING"
	case OrderProcessed:
		status = "PROCESSED"
	case OrderInvalid:
		status = "INVALID"
	}

	// Alias type is used to convert ProcessedAt to required format
	type Alias Order
	return json.Marshal(&struct {
		Status     string `json:"status"`
		UploadedAt string `json:"uploaded_at"`
		*Alias
	}{
		Alias:      (*Alias)(o),
		Status:     status,
		UploadedAt: o.UploadedAt.Format(time.RFC3339),
	})
}
