package models

import (
	"testing"

	"github.com/jon-kamis/klogger"
)

func TestValidateCanSaveBill(t *testing.T) {
	method := "Bill_test.TestValidateCanSaveBill"
	klogger.Enter(method)

	var b Bill

	err := b.ValidateCanSaveBill()

	if err == nil {
		t.Errorf("expected error to be thrown for empty bill but none was thrown")
	}

	b = Bill{
		UserID: 1,
		Amount: 1,
		Name:   "B1",
	}

	err = b.ValidateCanSaveBill()

	if err != nil {
		t.Errorf("unexpected error thrown for valid test case")
	}

	//Name is required
	b1 := b
	b1.Name = ""
	err = b1.ValidateCanSaveBill()

	if err == nil {
		t.Errorf("expected error to be thrown for empty bill name but none was thrown")
	}

	//Amount cannot be negative
	b1 = b
	b1.Amount = -1
	err = b1.ValidateCanSaveBill()

	if err == nil {
		t.Errorf("expected error to be thrown for negative bill amount but none was thrown")
	}

	//UserId is required
	b1 = b
	b1.UserID = 0
	err = b1.ValidateCanSaveBill()

	if err == nil {
		t.Errorf("expected error to be thrown for empty userId amount but none was thrown")
	}

	klogger.Exit(method)
}
