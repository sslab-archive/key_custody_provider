package entity

import "time"

type Authentication struct {
	ID           uint64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Payload      string
	AuthCode     string
	IsVerified   bool
}