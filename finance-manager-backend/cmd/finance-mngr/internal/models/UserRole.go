package models

import "time"

type UserRole struct {
	ID           int       `json:"id"`
	UserId       int       `json:"userId"`
	RoleId       int       `json:"roleId"`
	Code         string    `json:"code"`
	CreateDt     time.Time `json:"-"`
	LastUpdateDt time.Time `json:"-"`
}
