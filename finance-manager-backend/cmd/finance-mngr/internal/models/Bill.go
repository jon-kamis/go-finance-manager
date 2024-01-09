package models

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"time"
)

type Bill struct {
	ID           int       `json:"id"`
	UserID       int       `json:"userId"`
	Name         string    `json:"name"`
	Amount       float64   `json:"amount"`
	CreateDt     time.Time `json:"createDt"`
	LastUpdateDt time.Time `json:"lastUpdateDt"`
}

func (b *Bill) ValidateCanSaveBill() error {
	method := "Bill.ValidateCanSaveBill"
	fmlogger.Enter(method)

	if b.Name == "" {
		err := errors.New("name is required")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	if b.Amount < 0 {
		err := errors.New("amount cannot be negative")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	if b.UserID <= 0 {
		err := errors.New("userId is required")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}
