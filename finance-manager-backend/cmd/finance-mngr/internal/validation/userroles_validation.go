package validation

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
)

func (fmv *FinanceManagerValidator) UserRoleExistsAndBelongsToUser(roleId int, userId int) error {
	method := "userroles_validation.UserRoleExistsAndBelongsToUser"
	fmlogger.Enter(method)

	userRole, err := fmv.DB.GetUserRoleByRoleIDAndUserID(roleId, userId)

	if err != nil {
		fmlogger.ExitError(method, "database call returned with unexpected error", err)
		return err
	}

	//Check if UserRole is not null
	if userRole.ID > 0 {

		//Check if user Role belongs to specified user
		err = fmv.UserRoleBelongsToUser(userRole, userId)

		if err != nil {
			fmlogger.ExitError(method, "user role does not belong to user", err)
			return err
		}

	} else {
		err = errors.New("userRole does not exist")
		fmlogger.ExitError(method, "userRole does not exist", err)
		return err
	}

	fmlogger.Exit(method)
	return nil

}

func (fmv *FinanceManagerValidator) UserRoleBelongsToUser(userRole models.UserRole, userId int) error {
	method := "userroles_validation.UserRoleBelongsToUser"
	fmlogger.Enter(method)

	if userRole.ID == 0 || userRole.UserId == 0 || userId == 0 || userRole.UserId != userId {
		err := errors.New("user role does not belong to the user requesting it")
		fmlogger.ExitError(method, "user role does not belong to the user requesting it", err)
		return errors.New("forbidden")
	}

	fmlogger.Exit(method)
	return nil
}
