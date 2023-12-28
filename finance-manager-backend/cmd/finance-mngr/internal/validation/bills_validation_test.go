package validation

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"testing"
)

func TestBillBelongsToUser(t *testing.T) {
	method := "bills_validation_test.TestBillBelongsToUser"
	fmlogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Bill{
		ID:     1,
		UserID: 1,
	}

	err := v.BillBelongsToUser(l, 1)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error was thrown", err)
		t.Errorf("Unexpected error when validating Bill belongs to user %v\n", err)
	}

	fmlogger.Enter(method)
}

func TestBillBelongsToUser_nullBill(t *testing.T) {
	method := "bills_validation_test.TestBillBelongsToUser"
	fmlogger.Enter(method)

	v := FinanceManagerValidator{}

	err := v.BillBelongsToUser(models.Bill{}, 1)

	if err == nil {
		fmlogger.ExitError(method, "expected error but nothing was thrown", nil)
		t.Errorf("expected error but nothing was thrown\n")
	}

	fmlogger.Enter(method)
}

func TestBillBelongsToUser_nullUser(t *testing.T) {
	method := "bills_validation_test.TestBillBelongsToUser"
	fmlogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Bill{
		ID:     1,
		UserID: 1,
	}

	err := v.BillBelongsToUser(l, 0)

	if err == nil {
		fmlogger.ExitError(method, "expected error but nothing was thrown", nil)
		t.Errorf("expected error but nothing was thrown\n")
	}

	fmlogger.Enter(method)
}

func TestBillBelongsToUser_doesNotBelongToUser(t *testing.T) {
	method := "bills_validation_test.TestBillBelongsToUser"
	fmlogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Bill{
		ID:     1,
		UserID: 1,
	}

	err := v.BillBelongsToUser(l, 2)

	if err == nil {
		fmlogger.ExitError(method, "expected error but nothing was thrown", nil)
		t.Errorf("expected error but nothing was thrown\n")
	}

	fmlogger.Enter(method)
}
