package models

import (
	"database/sql"
	"errors"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"time"
)

// User struct contains a data link between user's and their stocks
type UserStock struct {
	ID           int          `json:"id"`
	UserId       int          `json:"userId" gorm:"column:user_id"`
	Ticker       string       `json:"ticker"`
	Quantity     float64      `json:"quantity"`
	EffectiveDt  time.Time    `json:"effectiveDt" gorm:"column:effective_dt"`
	ExpirationDt sql.NullTime `json:"expirationDt" gorm:"column:expiration_dt"`
	CreateDt     time.Time    `json:"createDt"`
	LastUpdateDt time.Time    `json:"lastUpdateDt"`
}

func (u *UserStock) ValidateCanSaveUserStock() error {
	method := "UserStock.ValidateCanSaveUserStock"
	fmlogger.Enter(method)

	var err error

	if u.UserId <= 0 {
		err = errors.New("userId is required")
		fmlogger.ExitError(method, "userId is required", err)
		return err
	}

	if u.Ticker == "" {
		err = errors.New("ticker is required")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	if u.Quantity < 0 {
		err = errors.New("quantity must be at least 0")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}
