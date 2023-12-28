package validation

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"testing"
)

func TestLoanBelongsToUser(t *testing.T) {
	method := "loans_validation_test.TestLoanBelongsToUser"
	fmlogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Loan{
		ID:     1,
		UserID: 1,
	}

	err := v.LoanBelongsToUser(l, 1)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error was thrown", err)
		t.Errorf("Unexpected error when validating Loan belongs to user %v\n", err)
	}

	fmlogger.Enter(method)
}

func TestLoanBelongsToUser_nullLoan(t *testing.T) {
	method := "loans_validation_test.TestLoanBelongsToUser"
	fmlogger.Enter(method)

	v := FinanceManagerValidator{}

	err := v.LoanBelongsToUser(models.Loan{}, 1)

	if err == nil {
		fmlogger.ExitError(method, "expected error but nothing was thrown", nil)
		t.Errorf("expected error but nothing was thrown\n")
	}

	fmlogger.Enter(method)
}

func TestLoanBelongsToUser_nullUser(t *testing.T) {
	method := "loans_validation_test.TestLoanBelongsToUser"
	fmlogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Loan{
		ID:     1,
		UserID: 1,
	}

	err := v.LoanBelongsToUser(l, 0)

	if err == nil {
		fmlogger.ExitError(method, "expected error but nothing was thrown", nil)
		t.Errorf("expected error but nothing was thrown\n")
	}

	fmlogger.Enter(method)
}

func TestLoanBelongsToUser_doesNotBelongToUser(t *testing.T) {
	method := "loans_validation_test.TestLoanBelongsToUser"
	fmlogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Loan{
		ID:     1,
		UserID: 1,
	}

	err := v.LoanBelongsToUser(l, 2)

	if err == nil {
		fmlogger.ExitError(method, "expected error but nothing was thrown", nil)
		t.Errorf("expected error but nothing was thrown\n")
	}

	fmlogger.Enter(method)
}
