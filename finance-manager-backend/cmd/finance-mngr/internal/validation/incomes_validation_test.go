package validation

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"testing"
)

func TestIncomeBelongsToUser(t *testing.T) {
	method := "incomes_validation_test.TestIncomeBelongsToUser"
	fmlogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Income{
		ID:     1,
		UserID: 1,
	}

	err := v.IncomeBelongsToUser(l, 1)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error was thrown", err)
		t.Errorf("Unexpected error when validating Income belongs to user %v\n", err)
	}

	fmlogger.Enter(method)
}

func TestIncomeBelongsToUser_nullIncome(t *testing.T) {
	method := "incomes_validation_test.TestIncomeBelongsToUser"
	fmlogger.Enter(method)

	v := FinanceManagerValidator{}

	err := v.IncomeBelongsToUser(models.Income{}, 1)

	if err == nil {
		fmlogger.ExitError(method, "expected error but nothing was thrown", nil)
		t.Errorf("expected error but nothing was thrown\n")
	}

	fmlogger.Enter(method)
}

func TestIncomeBelongsToUser_nullUser(t *testing.T) {
	method := "incomes_validation_test.TestIncomeBelongsToUser"
	fmlogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Income{
		ID:     1,
		UserID: 1,
	}

	err := v.IncomeBelongsToUser(l, 0)

	if err == nil {
		fmlogger.ExitError(method, "expected error but nothing was thrown", nil)
		t.Errorf("expected error but nothing was thrown\n")
	}

	fmlogger.Enter(method)
}

func TestIncomeBelongsToUser_doesNotBelongToUser(t *testing.T) {
	method := "incomes_validation_test.TestIncomeBelongsToUser"
	fmlogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Income{
		ID:     1,
		UserID: 1,
	}

	err := v.IncomeBelongsToUser(l, 2)

	if err == nil {
		fmlogger.ExitError(method, "expected error but nothing was thrown", nil)
		t.Errorf("expected error but nothing was thrown\n")
	}

	fmlogger.Enter(method)
}
