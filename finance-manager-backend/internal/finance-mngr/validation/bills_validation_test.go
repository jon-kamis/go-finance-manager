package validation

import (
	"finance-manager-backend/internal/finance-mngr/models"
	"testing"

	"github.com/jon-kamis/klogger"
)

func TestBillBelongsToUser(t *testing.T) {
	method := "bills_validation_test.TestBillBelongsToUser"
	klogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Bill{
		ID:     1,
		UserID: 1,
	}

	err := v.BillBelongsToUser(l, 1)

	if err != nil {
		t.Errorf("Unexpected error when validating Bill belongs to user %v\n", err)
	}

	klogger.Enter(method)
}

func TestBillBelongsToUser_nullBill(t *testing.T) {
	method := "bills_validation_test.TestBillBelongsToUser"
	klogger.Enter(method)

	v := FinanceManagerValidator{}

	err := v.BillBelongsToUser(models.Bill{}, 1)

	if err == nil {
		t.Errorf("expected error but nothing was thrown\n")
	}

	klogger.Enter(method)
}

func TestBillBelongsToUser_nullUser(t *testing.T) {
	method := "bills_validation_test.TestBillBelongsToUser"
	klogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Bill{
		ID:     1,
		UserID: 1,
	}

	err := v.BillBelongsToUser(l, 0)

	if err == nil {
		t.Errorf("expected error but nothing was thrown\n")
	}

	klogger.Enter(method)
}

func TestBillBelongsToUser_doesNotBelongToUser(t *testing.T) {
	method := "bills_validation_test.TestBillBelongsToUser"
	klogger.Enter(method)

	v := FinanceManagerValidator{}

	l := models.Bill{
		ID:     1,
		UserID: 1,
	}

	err := v.BillBelongsToUser(l, 2)

	if err == nil {
		t.Errorf("expected error but nothing was thrown\n")
	}

	klogger.Enter(method)
}
