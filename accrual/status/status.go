package status

import "strings"

type AccrualStatus string

const (
	Unknown    AccrualStatus = ""
	Registered AccrualStatus = "registered"
	Invalid    AccrualStatus = "invalid"
	Processing AccrualStatus = "processing"
	Processed  AccrualStatus = "processed"
)

func New(status string) AccrualStatus {
	switch strings.ToLower(status) {
	case "registered":
		return Registered
	case "invalid":
		return Invalid
	case "processing":
		return Processing
	case "processed":
		return Processed
	}
	return Unknown
}
