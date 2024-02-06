package validation

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"

	"github.com/jon-kamis/klogger"
)

func (fmv *FinanceManagerValidator) UserRoleExistsAndBelongsToUser(roleId int, userId int) error {
	method := "userroles_validation.UserRoleExistsAndBelongsToUser"
	klogger.Enter(method)

	userRole, err := fmv.DB.GetUserRoleByRoleIDAndUserID(roleId, userId)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	//Check if UserRole is not null
	if userRole.ID > 0 {

		//Check if user Role belongs to specified user
		err = fmv.UserRoleBelongsToUser(userRole, userId)

		if err != nil {
			klogger.ExitError(method, "user role does not belong to user")
			return err
		}

	} else {
		err = errors.New("userRole does not exist")
		klogger.ExitError(method, "userRole does not exist")
		return err
	}

	klogger.Exit(method)
	return nil

}

func (fmv *FinanceManagerValidator) UserRoleBelongsToUser(userRole models.UserRole, userId int) error {
	method := "userroles_validation.UserRoleBelongsToUser"
	klogger.Enter(method)

	if userRole.ID == 0 || userRole.UserId == 0 || userId == 0 || userRole.UserId != userId {
		klogger.ExitError(method, "user role does not belong to the user requesting it")
		return errors.New("forbidden")
	}

	klogger.Exit(method)
	return nil
}
