package validation

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"slices"
)

func (fmv *FinanceManagerValidator) IsValidToEnterNewUser(user models.User) error {
	method := "validate_users.isValidToEnterNewUser"
	fmt.Printf("[ENTER %s]\n", method)

	//validate required fields are present
	if user.Username == "" {
		errMsg := "username is required"
		fmt.Printf("[%s] %s", method, errMsg)
		fmt.Printf("[EXIT %s]\n", method)
		return errors.New(errMsg)
	}

	if user.Email == "" {
		errMsg := "email is required"
		fmt.Printf("[%s] %s", method, errMsg)
		fmt.Printf("[EXIT %s]\n", method)
		return errors.New(errMsg)
	}

	if user.Password == "" {
		errMsg := "password is required"
		fmt.Printf("[%s] %s", method, errMsg)
		fmt.Printf("[EXIT %s]\n", method)
		return errors.New(errMsg)
	}

	if user.FirstName == "" {
		errMsg := "first name is required"
		fmt.Printf("[%s] %s", method, errMsg)
		fmt.Printf("[EXIT %s]\n", method)
		return errors.New(errMsg)
	}

	if user.LastName == "" {
		errMsg := "last name is required"
		fmt.Printf("[%s] %s", method, errMsg)
		fmt.Printf("[EXIT %s]\n", method)
		return errors.New(errMsg)
	}

	// Check if username or email are valid
	_, err := fmv.DB.GetUserByUsernameOrEmail(user.Username, user.Email)
	if err == nil {
		errMsg := "username or email is not available"
		fmt.Printf("[%s] %s\n", method, errMsg)
		fmt.Printf("[EXIT %s]\n", method)
		return errors.New(errMsg)
	}

	fmt.Printf("[EXIT %s]\n", method)
	return nil
}

func (fmv *FinanceManagerValidator) IsValidToViewOtherUserData(loggedInUserId int) (bool, error) {
	method := "validate_users.isValidToViewOtherUserData"
	fmt.Printf("[ENTER %s]\n", method)

	valid, err := fmv.CheckIfUserHasRole(loggedInUserId, "admin")

	if err != nil {
		fmt.Printf("[%s] unexpected error returned while checking if user possesed desired role\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		return false, err
	}

	fmt.Printf("[EXIT %s]\n", method)
	return valid, nil
}

func (fmv *FinanceManagerValidator) CheckIfUserHasRole(id int, desiredRole string) (bool, error) {
	method := "utils.CheckIfUserHasRole"
	fmt.Printf("[ENTER %s]\n", method)

	userRoles, err := fmv.DB.GetUserRoles(id)

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching user roles\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		return false, err
	}

	var userRoleCodes []string
	for _, userRole := range userRoles {
		userRoleCodes = append(userRoleCodes, userRole.Code)
	}

	fmt.Printf("[EXIT %s]\n", method)
	return slices.Contains(userRoleCodes, desiredRole), nil
}
