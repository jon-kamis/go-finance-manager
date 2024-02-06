package validation

import (
	"finance-manager-backend/internal/finance-mngr/models"
	"testing"

	"github.com/jon-kamis/klogger"
)

func TestLoanBelongsToUser(t *testing.T) {
	method := "loans_validation_test.TestLoanBelongsToUser"
	klogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Loan{
		ID:     1,
		UserID: 1,
	}

	err := v.LoanBelongsToUser(l, 1)

	if err != nil {
		t.Errorf("Unexpected error when validating Loan belongs to user %v\n", err)
	}

	klogger.Enter(method)
}

func TestLoanBelongsToUser_nullLoan(t *testing.T) {
	method := "loans_validation_test.TestLoanBelongsToUser"
	klogger.Enter(method)

	v := FinanceManagerValidator{}

	err := v.LoanBelongsToUser(models.Loan{}, 1)

	if err == nil {
		t.Errorf("expected error but nothing was thrown\n")
	}

	klogger.Enter(method)
}

func TestLoanBelongsToUser_nullUser(t *testing.T) {
	method := "loans_validation_test.TestLoanBelongsToUser"
	klogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Loan{
		ID:     1,
		UserID: 1,
	}

	err := v.LoanBelongsToUser(l, 0)

	if err == nil {
		t.Errorf("expected error but nothing was thrown\n")
	}

	klogger.Enter(method)
}

func TestLoanBelongsToUser_doesNotBelongToUser(t *testing.T) {
	method := "loans_validation_test.TestLoanBelongsToUser"
	klogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Loan{
		ID:     1,
		UserID: 1,
	}

	err := v.LoanBelongsToUser(l, 2)

	if err == nil {
		t.Errorf("expected error but nothing was thrown\n")
	}

	klogger.Enter(method)
}
