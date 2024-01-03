package handlers

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/constants"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"net/http"
	"strconv"
)

func (fmh *FinanceManagerHandler) GetAndValidateUserId(idStr string, w http.ResponseWriter, r *http.Request) (int, error) {
	method := "handler_utils.getAndvalidateUserId"
	fmlogger.Enter(method)

	// Get loggedIn userId
	loggedInUserId, err := fmh.Auth.GetLoggedInUserId(w, r)
	var id int

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when retrieving logged in userId from HTTP Header", err)
		return -1, errors.New("failed to retrieve logged in userId")
	}

	if idStr == "me" {
		id = loggedInUserId
	} else {

		id, err = strconv.Atoi(idStr)
		if err != nil {
			fmlogger.ExitError(method, constants.FailedToParseIdError, err)
			return -1, err
		}

		canViewOtherUserData, err := fmh.CanViewOtherUserData(w, r)
		if err != nil {
			fmlogger.ExitError(method, err.Error(), err)
			return -1, err
		}

		if id != loggedInUserId && !canViewOtherUserData {
			err = errors.New(constants.UserForbiddenToViewOtherUserDataError)
			fmlogger.ExitError(method, constants.UserForbiddenToViewOtherUserDataError, err)
			return -1, err
		}
	}

	fmlogger.Exit(method)
	return id, nil
}

func (fmh *FinanceManagerHandler) CanViewOtherUserData(w http.ResponseWriter, r *http.Request) (bool, error) {
	method := "handler_utils.CanViewOtherUserData"
	fmlogger.Enter(method)

	loggedInUserId, err := fmh.Auth.GetLoggedInUserId(w, r)

	if err != nil {
		fmlogger.ExitError(method, constants.FailedToReadUserIdFromAuthHeaderError, err)
		return false, errors.New(constants.FailedToReadUserIdFromAuthHeaderError)
	}

	//Fetch User Roles to see if this user is an Administrator
	isValid, err := fmh.Validator.CheckIfUserHasRole(loggedInUserId, "admin")

	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		return false, err
	}

	fmlogger.Exit(method)
	return isValid, nil
}
