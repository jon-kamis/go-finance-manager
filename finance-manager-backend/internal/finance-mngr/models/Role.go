package models

import "time"

type Role struct {
	ID           int       `json:"id"`
	Code         string    `json:"code"`
	CreateDt     time.Time `json:"-"`
	LastUpdateDt time.Time `json:"-"`
}
