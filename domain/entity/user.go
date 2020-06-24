package entity

import "time"

type User struct {
	ID           uint64
	PublicKey    string
	PartialKey   string
	RegisteredTx string

	CreatedAt    time.Time
	UpdatedAt    time.Time
}
