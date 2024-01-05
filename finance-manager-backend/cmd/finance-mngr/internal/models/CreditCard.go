package models

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/constants"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"time"
)

type CreditCard struct {
	ID                   int       `json:"id"`
	UserID               int       `json:"userId" gorm:"user_id"`
	Name                 string    `json:"name"`
	Balance              float64   `json:"balance"`
	Limit                float64   `json:"limit" gorm:"column:credit_limit"`
	APR                  float64   `json:"apr"`
	MinPayment           float64   `json:"minPayment" gorm:"column:min_pay"`
	MinPaymentPercentage float64   `json:"minPaymentPercentage" gorm:"column:min_pay_percentage"`
	CreateDt             time.Time `json:"createDt" gorm:"column:create_dt"`
	LastUpdateDt         time.Time `json:"lastUpdateDt" gorm:"column:last_update_dt"`
}

func (cc *CreditCard) ValidateCanSaveCreditCard() error {
	method := "creditcard.ValidateCanSaveCreditCard"
	fmlogger.Enter(method)

	if cc.ID == 0 || cc.UserID == 0 || cc.Name == "" || cc.MinPayment == 0 || cc.MinPaymentPercentage == 0 {
		err := errors.New(constants.InvalidCreditCardError)
		fmlogger.ExitError(method, "one or more required fields is blank", err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}
