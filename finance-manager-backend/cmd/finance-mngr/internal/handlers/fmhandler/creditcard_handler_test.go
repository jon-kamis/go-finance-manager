package fmhandler

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"testing"
)

func TestGetAllUserCreditCards(t *testing.T) {
	method := "creditcard_handler_test.TestGetAllUserCreditCards"
	fmlogger.Enter(method)

	fmlogger.Exit(method)
}

func TestSaveCreditCard(t *testing.T) {
	method := "creditcard_handler_test.TestSaveCreditCard"
	fmlogger.Enter(method)
	/*
		cc := models.CreditCard{
			UserID: 1,
			Name: "TestSaveCreditCard",
			Balance: 1000.0,
			Limit: 20000.0,
			APR: 26.2,
			MinPayment: 35.00,
			MinPaymentPercentage: 10,
		}

		writer := MakeRequest(http.MethodGet, "")
	*/
	fmlogger.Exit(method)
}
