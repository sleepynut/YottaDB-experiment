package model

import "time"

type Contact struct {
	ContactID     string
	Type          string
	IsPrimary     bool
	UserID        string
	AddressInfo   string
	Province      string
	Tel           string
	CreatedDt     time.Time
	LastUpdatedDt time.Time
}
