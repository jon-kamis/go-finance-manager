package validation

import (
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/test"
	"testing"
)

func TestCreditCardBelongsToUser(t *testing.T) {
	method := "creditcards_validation_test.TestCreditCardBelongsToUser"
	fmlogger.Enter(method)

	userId := test.TestingAdmin.ID

	cc := models.CreditCard{
		ID:                   1,
		UserID:               userId,
		Name:                 "Testing Card",
		Balance:              1000.0,
		APR:                  26.0,
		MinPayment:           35.00,
		MinPaymentPercentage: 10,
	}

	//cc is not initialized
	err := fmv.CreditCardBelongsToUser(models.CreditCard{}, userId)

	if err == nil {
		t.Errorf("expected error to be thrown for uninitialized credit card but none was thrown")
	}

	err = fmv.CreditCardBelongsToUser(cc, 0)

	if err == nil {
		t.Errorf("expected error to be thrown for invalid userId but none was thrown")
	}

	err = fmv.CreditCardBelongsToUser(cc, 2)

	if err == nil {
		t.Errorf("expected error to be thrown for credit card does not belong to user but none was thrown")
	}

	err = fmv.CreditCardBelongsToUser(cc, userId)

	if err != nil {
		t.Errorf("unexpected error thrown for valid test case")
	}

	fmlogger.Exit(method)

}
