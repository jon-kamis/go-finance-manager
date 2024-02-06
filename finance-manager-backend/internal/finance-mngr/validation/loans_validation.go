package validation

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/models"

	"github.com/jon-kamis/klogger"
)

func (fmv *FinanceManagerValidator) LoanBelongsToUser(loan models.Loan, userId int) error {
	method := "loans_validation.isValidToReturnLoanToUser"
	klogger.Enter(method)

	if loan.ID == 0 || loan.UserID == 0 || userId == 0 || loan.UserID != userId {
		klogger.ExitError(method, "loan does not belong to user")
		return errors.New("forbidden")
	}

	klogger.Exit(method)
	return nil
}
