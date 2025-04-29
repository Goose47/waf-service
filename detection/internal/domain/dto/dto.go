// Package dto contains DTOs used in service.
package dto

import "time"

// Client represents a single client and contains its requests' timestamps.
type Client struct {
	IP           string      `json:"ip"`
	Fingerprints []time.Time `json:"fingerprints"`
	IsSuspicious bool        `json:"is_suspicious"`
}
