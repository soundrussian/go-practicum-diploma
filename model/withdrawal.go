package model

import (
	"encoding/json"
	"time"
)

type Withdrawal struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"-"`
}

// MarshalJSON is used to convert ProcessedAt to required RFC3339 format on marshal
func (d *Withdrawal) MarshalJSON() ([]byte, error) {
	// Alias type is used to convert ProcessedAt to required format
	type Alias Withdrawal
	return json.Marshal(&struct {
		*Alias
		Timestamp string `json:"processed_at"`
	}{
		Alias:     (*Alias)(d),
		Timestamp: d.ProcessedAt.Format(time.RFC3339),
	})
}
