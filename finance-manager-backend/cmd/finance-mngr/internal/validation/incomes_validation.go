package validation

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
)

func (fmv *FinanceManagerValidator) IncomeBelongsToUser(income models.Income, userId int) error {
	method := "loans_validation.isValidToReturnLoanToUser"
	fmt.Printf("[ENTER %s]\n", method)

	if income.ID == 0 || income.UserID == 0 || userId == 0 || income.UserID != userId {
		fmt.Printf("[%s] loan does not belong to user requesting it!\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		return errors.New("forbidden")
	}

	fmt.Printf("[EXIT %s]\n", method)
	return nil
}
