package model

import "time"

// WorldConfig represents a single key-value configuration entry.
type WorldConfig struct {
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}
