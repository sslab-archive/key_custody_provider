package entity

import "time"

type User struct {
	ID              uint64
	PublicKey       string
	PartialKey      string
	PartialKeyIndex uint64
	Payload         string
	RegisteredTx    string

	CreatedAt time.Time
	UpdatedAt time.Time
}
