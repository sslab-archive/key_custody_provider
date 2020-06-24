package entity

import "time"

type Credential struct {
	ID         uint64
	Payload    string
	PartialKey string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
