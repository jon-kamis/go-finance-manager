package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"net/http"
	"strconv"

	"github.com/jon-kamis/klogger"
)

func (fmh *FinanceManagerHandler) GetAndValidateUserId(idStr string, w http.ResponseWriter, r *http.Request) (int, error) {
	method := "handler_utils.getAndvalidateUserId"
	klogger.Enter(method)

	// Get loggedIn userId
	loggedInUserId, err := fmh.Auth.GetLoggedInUserId(w, r)
	var id int

	if err != nil {
		klogger.ExitError(method, constants.FailedToReadUserIdFromAuthHeaderError, err)
		return -1, errors.New("failed to retrieve logged in userId")
	}

	if idStr == "me" {
		id = loggedInUserId
	} else {

		id, err = strconv.Atoi(idStr)
		if err != nil {
			klogger.ExitError(method, constants.FailedToParseIdError, err)
			return -1, err
		}

		canViewOtherUserData, err := fmh.CanViewOtherUserData(w, r)
		if err != nil {
			klogger.ExitError(method, err.Error(), err)
			return -1, err
		}

		if id != loggedInUserId && !canViewOtherUserData {
			err = errors.New(constants.UserForbiddenToViewOtherUserDataError)
			klogger.ExitError(method, constants.UserForbiddenToViewOtherUserDataError, err)
			return -1, err
		}
	}

	klogger.Exit(method)
	return id, nil
}

func (fmh *FinanceManagerHandler) CanViewOtherUserData(w http.ResponseWriter, r *http.Request) (bool, error) {
	method := "handler_utils.CanViewOtherUserData"
	klogger.Enter(method)

	loggedInUserId, err := fmh.Auth.GetLoggedInUserId(w, r)

	if err != nil {
		klogger.ExitError(method, constants.FailedToReadUserIdFromAuthHeaderError, err)
		return false, errors.New(constants.FailedToReadUserIdFromAuthHeaderError)
	}

	//Fetch User Roles to see if this user is an Administrator
	isValid, err := fmh.Validator.CheckIfUserHasRole(loggedInUserId, "admin")

	if err != nil {
		klogger.ExitError(method, err.Error(), err)
		return false, err
	}

	klogger.Exit(method)
	return isValid, nil
}
