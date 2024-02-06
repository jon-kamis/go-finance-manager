package validation

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/models"

	"github.com/jon-kamis/klogger"
)

func (fmv *FinanceManagerValidator) BillBelongsToUser(bill models.Bill, userId int) error {
	method := "loans_validation.isValidToReturnLoanToUser"
	klogger.Enter(method)

	if bill.ID == 0 || bill.UserID == 0 || userId == 0 || bill.UserID != userId {
		err := errors.New("forbidden")
		klogger.ExitError(method, "bill does not belong to logged in user")
		return err
	}

	klogger.Exit(method)
	return nil
}
