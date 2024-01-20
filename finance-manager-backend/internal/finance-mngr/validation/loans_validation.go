package validation

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/models"
	"fmt"
)

func (fmv *FinanceManagerValidator) LoanBelongsToUser(loan models.Loan, userId int) error {
	method := "loans_validation.isValidToReturnLoanToUser"
	fmt.Printf("[ENTER %s]\n", method)

	if loan.ID == 0 || loan.UserID == 0 || userId == 0 || loan.UserID != userId {
		fmt.Printf("[%s] loan does not belong to user requesting it!\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		return errors.New("forbidden")
	}

	fmt.Printf("[EXIT %s]\n", method)
	return nil
}
