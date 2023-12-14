package models

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"fmt"
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
	fmt.Printf("[ENTER %s]\n", method)

	if b.Name == "" {
		returnError("cannot save bill without a name", method)
	}

	if b.Amount < 0 {
		returnError("Gross Pay is required", method)
	}

	if b.UserID <= 0 {
		returnError("UserId is required", method)
	}

	fmlogger.Exit(method)
	return nil
}
