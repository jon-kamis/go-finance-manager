package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jon-kamis/klogger"
)

// UserStock struct contains a data link between user's and their stocks
type UserStock struct {
	ID           int          `json:"id"`
	UserId       int          `json:"userId" gorm:"column:user_id"`
	Ticker       string       `json:"ticker"`
	Quantity     float64      `json:"quantity"`
	Type         string       `json:"type"`
	EffectiveDt  time.Time    `json:"effectiveDt" gorm:"column:effective_dt"`
	ExpirationDt sql.NullTime `json:"expirationDt" gorm:"column:expiration_dt"`
	CreateDt     time.Time    `json:"createDt"`
	LastUpdateDt time.Time    `json:"lastUpdateDt"`
}

func (u *UserStock) ValidateCanSaveUserStock() error {
	method := "UserStock.ValidateCanSaveUserStock"
	klogger.Enter(method)

	var err error

	if u.UserId <= 0 {
		err = errors.New("userId is required")
		klogger.ExitError(method, "userId is required")
		return err
	}

	if u.Ticker == "" {
		err = errors.New("ticker is required")
		klogger.ExitError(method, err.Error())
		return err
	}

	if u.Quantity < 0 {
		err = errors.New("quantity must be at least 0")
		klogger.ExitError(method, err.Error())
		return err
	}

	klogger.Exit(method)
	return nil
}
