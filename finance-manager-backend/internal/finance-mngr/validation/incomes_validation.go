package validation

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/models"

	"github.com/jon-kamis/klogger"
)

func (fmv *FinanceManagerValidator) IncomeBelongsToUser(income models.Income, userId int) error {
	method := "incomes_validation.IncomeBelongsToUser"
	klogger.Enter(method)

	if income.ID == 0 || income.UserID == 0 || userId == 0 || income.UserID != userId {
		klogger.ExitError(method, "income does not belong to user")
		return errors.New("forbidden")
	}

	klogger.Exit(method)
	return nil
}
