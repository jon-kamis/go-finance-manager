package validation

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/models"

	"github.com/jon-kamis/klogger"
)

func (fmv *FinanceManagerValidator) CreditCardBelongsToUser(cc models.CreditCard, userId int) error {
	method := "creditcards_validation.CreditCardBelongsToUser"
	klogger.Enter(method)

	if cc.ID == 0 || cc.UserID == 0 || userId == 0 || cc.UserID != userId {
		err := errors.New("forbidden")
		klogger.ExitError(method, "credit card does not belong to logged in user")
		return err
	}

	klogger.Exit(method)
	return nil
}
