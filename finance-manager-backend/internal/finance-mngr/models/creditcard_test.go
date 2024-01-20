package models

import (
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestCalcPayment(t *testing.T) {

	cc := CreditCard{
		Balance:              100,
		MinPayment:           35,
		MinPaymentPercentage: 10,
	}

	//Min Payment should be the greater value between (balance * minPaymentPercentage) and MinPayment

	//MinPayment is greater
	cc.CalcPayment()
	assert.Equal(t, cc.MinPayment, cc.Payment)

	cc.Balance = 1000
	cc.CalcPayment()
	assert.Equal(t, 100.0, cc.Payment)

}
