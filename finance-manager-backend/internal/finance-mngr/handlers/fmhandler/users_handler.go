package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jon-kamis/klogger"
)

// GetUserByID godoc
// @title		Get User by ID
// @version 	1.0.0
// @Tags 		Users
// @Summary 	Get User by ID
// @Description Returns a User by its ID
// @Param		userId path int true "ID of the user to fetch"
// @Produce 	json
// @Success 	200 {object} models.User
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId} [get]
func (fmh *FinanceManagerHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	method := "users_handler.GetUserByID"
	klogger.Enter(method)

	idStr := chi.URLParam(r, "userId")
	var id int

	loggedInUserId, err := fmh.Auth.GetLoggedInUserId(w, r)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching user"), http.StatusInternalServerError)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	if idStr == "me" {
		id = loggedInUserId
	} else {

		id, err = strconv.Atoi(idStr)
		if err != nil {
			fmh.JSONUtil.ErrorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
			klogger.ExitError(method, constants.ProcessIdError, err)
			return
		}

		if id != loggedInUserId {
			isValid, err := fmh.Validator.IsValidToViewOtherUserData(loggedInUserId)

			if err != nil {
				fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching user"), http.StatusInternalServerError)
				klogger.ExitError(method, constants.GenericUnexpectedErrorLog, err)

				return
			} else if !isValid {
				fmh.JSONUtil.ErrorJSON(w, errors.New("forbidden"), http.StatusForbidden)
				klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
				return
			}

		}
	}

	user, err := fmh.DB.GetUserByID(id)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.GenericServerError), http.StatusInternalServerError)
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return
	}

	if user.ID <= 0 {
		fmh.JSONUtil.ErrorJSON(w, errors.New("user not found"), http.StatusNotFound)
		klogger.ExitError(method, constants.EntityNotFoundError)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, user)
}

// DeleteUserById godoc
// @title		Delete User by ID
// @version 	1.0.0
// @Tags 		Users
// @Summary 	Delete User by ID
// @Description Deletes a User by its ID. Cascades to all objects owned by the user
// @Param		userId path int true "ID of the user to fetch"
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId} [delete]
func (fmh *FinanceManagerHandler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	method := "users_handler.DeleteUserById"
	klogger.Enter(method)

	idStr := chi.URLParam(r, "userId")
	var id int

	loggedInUserId, err := fmh.Auth.GetLoggedInUserId(w, r)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching user"), http.StatusInternalServerError)
		klogger.ExitError(method, constants.FailedToReadUserIdFromAuthHeaderError, err)
		return
	}

	if idStr == "me" {
		id = loggedInUserId
	} else {

		id, err = strconv.Atoi(idStr)
		if err != nil {
			fmh.JSONUtil.ErrorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
			klogger.ExitError(method, constants.ProcessIdError, err)
			return
		}

		if id != loggedInUserId {
			isValid, err := fmh.Validator.IsValidToViewOtherUserData(loggedInUserId)

			if err != nil {
				fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching user"), http.StatusInternalServerError)
				klogger.ExitError(method, constants.GenericUnexpectedErrorLog, err)

				return
			} else if !isValid {
				fmh.JSONUtil.ErrorJSON(w, errors.New("forbidden"), http.StatusForbidden)
				klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
				return
			}

		}
	}

	_, err = fmh.DB.GetUserByID(id)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New("user not found"), http.StatusNotFound)
		klogger.ExitError(method, constants.EntityNotFoundError)
		return
	}

	err = fmh.DB.DeleteBillsByUserID(id)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New("an unexpected error occured while attempting to delete the user"), http.StatusNotFound)
		klogger.ExitError(method, "failed to delete user bills:\n%v", err)
		return
	}

	err = fmh.DB.DeleteIncomesByUserID(id)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New("an unexpected error occured while attempting to delete the user"), http.StatusNotFound)
		klogger.ExitError(method, "failed to delete user incomes:\n%v", err)
		return
	}

	err = fmh.DB.DeleteLoansByUserID(id)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New("an unexpected error occured while attempting to delete the user"), http.StatusNotFound)
		klogger.ExitError(method, "failed to delete user loans:\n%v", err)
		return
	}

	err = fmh.DB.DeleteCreditCardsByUserID(id)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New("an unexpected error occured while attempting to delete the user"), http.StatusNotFound)
		klogger.ExitError(method, "failed to delete user credit cards:\n%v", err)
		return
	}

	err = fmh.DB.DeleteUserByID(id)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New("an unexpected error occured while attempting to delete the user"), http.StatusNotFound)
		klogger.ExitError(method, "failed to delete user:\n%v", err)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, "success")
}

// GetAllUsers godoc
// @title		Get All Users
// @version 	1.0.0
// @Tags 		Users
// @Summary 	Get All Users
// @Description Returns an array of User objects
// @Param		search query string false "Search for Users by first or last name"
// @Produce 	json
// @Success 	200 {array} models.User
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users [get]
func (fmh *FinanceManagerHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	method := "users_handler.GetAllUsers"
	klogger.Enter(method)

	search := r.URL.Query().Get("search")
	users, err := fmh.DB.GetAllUsers(search)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching users"), http.StatusInternalServerError)
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, users)
}
