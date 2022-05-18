package model

import "time"

type Transaction struct {
	UUID          string
	Status        string
	Src           string
	Dst           string
	Amount        float64
	CreatedDt     time.Time
	LastUpdatedDt time.Time
}
