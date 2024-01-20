package validation

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
)

func (fmv *FinanceManagerValidator) CreditCardBelongsToUser(cc models.CreditCard, userId int) error {
	method := "creditcards_validation.CreditCardBelongsToUser"
	fmlogger.Enter(method)

	if cc.ID == 0 || cc.UserID == 0 || userId == 0 || cc.UserID != userId {
		err := errors.New("forbidden")
		fmlogger.ExitError(method, "credit card does not belong to logged in user", err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}
