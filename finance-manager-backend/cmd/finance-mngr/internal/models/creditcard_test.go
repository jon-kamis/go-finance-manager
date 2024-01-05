package models

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"testing"
)

func TestValidateCanSaveCreditCard(t *testing.T) {
	method := "creditcard_test.TestValidateCanSaveCreditCard"
	fmlogger.Enter(method)

	//Valid Case
	cc := CreditCard{
		ID:                   1,
		UserID:               1,
		Name:                 "cc",
		MinPayment:           35.0,
		MinPaymentPercentage: 10,
	}

	err := cc.ValidateCanSaveCreditCard()
	if err != nil {
		t.Errorf("[%s] unexpected error was thrown during valid test case", err)
	}

	cc2 := CreditCard{}
	err = cc2.ValidateCanSaveCreditCard()

	if err == nil {
		t.Errorf("[%s] expected error for invalid credit card but none was thrown", err)
	}

	fmlogger.Exit(method)
}
