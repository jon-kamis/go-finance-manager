package validation

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"slices"

	"github.com/jon-kamis/klogger"
)

func (fmv *FinanceManagerValidator) IsValidToEnterNewUser(user models.User) error {
	method := "users_validation.isValidToEnterNewUser"
	klogger.Enter(method)

	//validate required fields are present
	if user.Username == "" {
		errMsg := "username is required"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	if user.Email == "" {
		errMsg := "email is required"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	if user.Password == "" {
		errMsg := "password is required"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	if user.FirstName == "" {
		errMsg := "first name is required"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	if user.LastName == "" {
		errMsg := "last name is required"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	// Check if username or email are valid
	u, err := fmv.DB.GetUserByUsernameOrEmail(user.Username, user.Email)

	if err != nil {
		klogger.ExitError(method, constants.GenericUnexpectedErrorLog, err)
		return err
	}

	if u.ID > 0 {
		err := errors.New(constants.UsernameOrEmailExistError)
		klogger.ExitError(method, constants.UsernameOrEmailExistError)
		return err
	}

	klogger.Exit(method)
	return nil
}

func (fmv *FinanceManagerValidator) IsValidToViewOtherUserData(loggedInUserId int) (bool, error) {
	method := "users_validation.isValidToViewOtherUserData"
	klogger.Enter(method)

	valid, err := fmv.CheckIfUserHasRole(loggedInUserId, "admin")

	if err != nil {
		klogger.ExitError(method, "unexpected error retruned while checking if user possesed desired role:\n%v", err)
		return false, err
	}

	klogger.Exit(method)
	return valid, nil
}

func (fmv *FinanceManagerValidator) IsValidToDeleteOtherUserData(loggedInUserId int) (bool, error) {
	method := "users_validation.IsValidToDeleteOtherUserData"
	klogger.Enter(method)

	valid, err := fmv.CheckIfUserHasRole(loggedInUserId, "admin")

	if err != nil {
		klogger.ExitError(method, "unexpected error retruned while checking if user possesed desired role:\n%v", err)
		return false, err
	}

	klogger.Exit(method)
	return valid, nil
}

func (fmv *FinanceManagerValidator) CheckIfUserHasRole(id int, desiredRole string) (bool, error) {
	method := "users_validation.CheckIfUserHasRole"
	klogger.Enter(method)

	userRoles, err := fmv.DB.GetAllUserRoles(id)

	if err != nil {
		klogger.ExitError(method, constants.GenericUnexpectedErrorLog, err)
		return false, err
	}

	var userRoleCodes []string
	for _, userRole := range userRoles {
		userRoleCodes = append(userRoleCodes, userRole.Code)
	}

	klogger.Exit(method)
	return slices.Contains(userRoleCodes, desiredRole), nil
}
