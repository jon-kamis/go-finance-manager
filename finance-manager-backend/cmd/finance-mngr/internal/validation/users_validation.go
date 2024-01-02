package validation

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"slices"
)

func (fmv *FinanceManagerValidator) IsValidToEnterNewUser(user models.User) error {
	method := "users_validation.isValidToEnterNewUser"
	fmlogger.Enter(method)

	//validate required fields are present
	if user.Username == "" {
		errMsg := "username is required"
		fmt.Printf("[%s] %s", method, errMsg)
		fmlogger.Exit(method)
		return errors.New(errMsg)
	}

	if user.Email == "" {
		errMsg := "email is required"
		fmt.Printf("[%s] %s", method, errMsg)
		fmlogger.Exit(method)
		return errors.New(errMsg)
	}

	if user.Password == "" {
		errMsg := "password is required"
		fmt.Printf("[%s] %s", method, errMsg)
		fmlogger.Exit(method)
		return errors.New(errMsg)
	}

	if user.FirstName == "" {
		errMsg := "first name is required"
		fmt.Printf("[%s] %s", method, errMsg)
		fmlogger.Exit(method)
		return errors.New(errMsg)
	}

	if user.LastName == "" {
		errMsg := "last name is required"
		fmt.Printf("[%s] %s", method, errMsg)
		fmlogger.Exit(method)
		return errors.New(errMsg)
	}

	// Check if username or email are valid
	_, err := fmv.DB.GetUserByUsernameOrEmail(user.Username, user.Email)
	if err == nil {
		errMsg := "username or email is not available"
		fmt.Printf("[%s] %s\n", method, errMsg)
		fmlogger.Exit(method)
		return errors.New(errMsg)
	}

	fmlogger.Exit(method)
	return nil
}

func (fmv *FinanceManagerValidator) IsValidToViewOtherUserData(loggedInUserId int) (bool, error) {
	method := "users_validation.isValidToViewOtherUserData"
	fmlogger.Enter(method)

	valid, err := fmv.CheckIfUserHasRole(loggedInUserId, "admin")

	if err != nil {
		fmlogger.ExitError(method, "unexpected error retruned while checking if user possesed desired role", err)
		return false, err
	}

	fmlogger.Exit(method)
	return valid, nil
}

func (fmv *FinanceManagerValidator) IsValidToDeleteOtherUserData(loggedInUserId int) (bool, error) {
	method := "users_validation.IsValidToDeleteOtherUserData"
	fmlogger.Enter(method)

	valid, err := fmv.CheckIfUserHasRole(loggedInUserId, "admin")

	if err != nil {
		fmlogger.ExitError(method, "unexpected error retruned while checking if user possesed desired role", err)
		return false, err
	}

	fmlogger.Exit(method)
	return valid, nil
}

func (fmv *FinanceManagerValidator) CheckIfUserHasRole(id int, desiredRole string) (bool, error) {
	method := "users_validation.CheckIfUserHasRole"
	fmlogger.Enter(method)

	userRoles, err := fmv.DB.GetAllUserRoles(id)

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching user roles\n", method)
		fmlogger.Exit(method)
		return false, err
	}

	var userRoleCodes []string
	for _, userRole := range userRoles {
		userRoleCodes = append(userRoleCodes, userRole.Code)
	}

	fmlogger.Exit(method)
	return slices.Contains(userRoleCodes, desiredRole), nil
}
