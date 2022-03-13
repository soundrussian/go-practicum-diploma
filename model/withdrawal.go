package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Withdrawal struct {
	Order       string
	Sum         decimal.Decimal
	ProcessedAt time.Time
}
