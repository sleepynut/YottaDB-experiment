package model

import "time"

type Transaction struct {
	UID           string
	Status        string
	CreatedDt     time.Time
	LastUpdatedDt time.Time
}
