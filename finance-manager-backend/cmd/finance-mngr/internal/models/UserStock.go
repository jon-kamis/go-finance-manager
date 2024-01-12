package models

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"time"
)

// User struct contains a data link between user's and their stocks
type UserStock struct {
	ID           int       `json:"id"`
	UserId       int       `json:"userId" gorm:"column:user_id"`
	Ticker       string    `json:"ticker"`
	Quantity     float64   `json:"quantity"`
	CreateDt     time.Time `json:"createDt"`
	LastUpdateDt time.Time `json:"lastUpdateDt"`
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
