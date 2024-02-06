package validation

import (
	"finance-manager-backend/internal/finance-mngr/models"
	"testing"

	"github.com/jon-kamis/klogger"
)

func TestIncomeBelongsToUser(t *testing.T) {
	method := "incomes_validation_test.TestIncomeBelongsToUser"
	klogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Income{
		ID:     1,
		UserID: 1,
	}

	err := v.IncomeBelongsToUser(l, 1)

	if err != nil {
		t.Errorf("Unexpected error when validating Income belongs to user %v\n", err)
	}

	klogger.Enter(method)
}

func TestIncomeBelongsToUser_nullIncome(t *testing.T) {
	method := "incomes_validation_test.TestIncomeBelongsToUser"
	klogger.Enter(method)

	v := FinanceManagerValidator{}

	err := v.IncomeBelongsToUser(models.Income{}, 1)

	if err == nil {
		t.Errorf("expected error but nothing was thrown\n")
	}

	klogger.Enter(method)
}

func TestIncomeBelongsToUser_nullUser(t *testing.T) {
	method := "incomes_validation_test.TestIncomeBelongsToUser"
	klogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Income{
		ID:     1,
		UserID: 1,
	}

	err := v.IncomeBelongsToUser(l, 0)

	if err == nil {
		t.Errorf("expected error but nothing was thrown\n")
	}

	klogger.Enter(method)
}

func TestIncomeBelongsToUser_doesNotBelongToUser(t *testing.T) {
	method := "incomes_validation_test.TestIncomeBelongsToUser"
	klogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Income{
		ID:     1,
		UserID: 1,
	}

	err := v.IncomeBelongsToUser(l, 2)

	if err == nil {
		t.Errorf("expected error but nothing was thrown\n")
	}

	klogger.Enter(method)
}
