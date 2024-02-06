package models

import (
	"errors"
	"time"

	"github.com/jon-kamis/klogger"
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
	klogger.Enter(method)

	if b.Name == "" {
		err := errors.New("name is required")
		klogger.ExitError(method, err.Error())
		return err
	}

	if b.Amount < 0 {
		err := errors.New("amount cannot be negative")
		klogger.ExitError(method, err.Error())
		return err
	}

	if b.UserID <= 0 {
		err := errors.New("userId is required")
		klogger.ExitError(method, err.Error())
		return err
	}

	klogger.Exit(method)
	return nil
}
