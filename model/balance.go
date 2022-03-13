package model

import "github.com/shopspring/decimal"

type UserBalance struct {
	Current   decimal.Decimal
	Withdrawn decimal.Decimal
}
