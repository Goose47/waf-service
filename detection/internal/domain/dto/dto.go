// Package dto contains DTOs used in service.
package dto

import "time"

type Client struct {
	IP           string      `json:"ip"`
	IsSuspicious bool        `json:"is_suspicious"`
	Fingerprints []time.Time `json:"fingerprints"`
}
