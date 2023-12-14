package validation

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
)

func (fmv *FinanceManagerValidator) BillBelongsToUser(bill models.Bill, userId int) error {
	method := "loans_validation.isValidToReturnLoanToUser"
	fmlogger.Enter(method)

	if bill.ID == 0 || bill.UserID == 0 || userId == 0 || bill.UserID != userId {
		err := errors.New("forbidden")
		fmlogger.ExitError(method, "bill does not belong to logged in user", err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}
