package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	fmlogger.Enter(method)

	idStr := chi.URLParam(r, "userId")
	var id int

	loggedInUserId, err := fmh.Auth.GetLoggedInUserId(w, r)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching logged in user", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching user"), http.StatusInternalServerError)
		return
	}

	if idStr == "me" {
		id = loggedInUserId
	} else {

		id, err = strconv.Atoi(idStr)
		if err != nil {
			fmlogger.ExitError(method, "the supplied id was invalid", err)
			fmh.JSONUtil.ErrorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
			return
		}

		if id != loggedInUserId {
			isValid, err := fmh.Validator.IsValidToViewOtherUserData(loggedInUserId)

			if err != nil {
				fmlogger.ExitError(method, "unexpected error occured when fetching the selected user", err)
				fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching user"), http.StatusInternalServerError)
				return
			} else if !isValid {
				fmlogger.ExitError(method, "user is forbidden to view other user data", err)
				fmh.JSONUtil.ErrorJSON(w, errors.New("forbidden"), http.StatusForbidden)
				return
			}

		}
	}

	user, err := fmh.DB.GetUserByID(id)

	if err != nil {
		fmlogger.ExitError(method, "requested user was not found", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("user not found"), http.StatusNotFound)
		return
	}

	fmlogger.Exit(method)
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
	fmlogger.Enter(method)

	idStr := chi.URLParam(r, "userId")
	var id int

	loggedInUserId, err := fmh.Auth.GetLoggedInUserId(w, r)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching logged in user", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching user"), http.StatusInternalServerError)
		return
	}

	if idStr == "me" {
		id = loggedInUserId
	} else {

		id, err = strconv.Atoi(idStr)
		if err != nil {
			fmlogger.ExitError(method, "the supplied id was invalid", err)
			fmh.JSONUtil.ErrorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
			return
		}

		if id != loggedInUserId {
			isValid, err := fmh.Validator.IsValidToDeleteOtherUserData(loggedInUserId)

			if err != nil {
				fmlogger.ExitError(method, "unexpected error occured when fetching the selected user", err)
				fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching user"), http.StatusInternalServerError)
				return
			} else if !isValid {
				fmlogger.ExitError(method, "user is forbidden to delete other user data", err)
				fmh.JSONUtil.ErrorJSON(w, errors.New("forbidden"), http.StatusForbidden)
				return
			}

		}
	}

	_, err = fmh.DB.GetUserByID(id)

	if err != nil {
		fmlogger.ExitError(method, "requested user was not found", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("user not found"), http.StatusNotFound)
		return
	}

	err = fmh.DB.DeleteBillsByUserID(id)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error when deleting bills by userId", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("an unexpected error occured while attempting to delete the user"), http.StatusNotFound)
		return
	}

	err = fmh.DB.DeleteIncomesByUserID(id)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error when deleting incomes by userId", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("an unexpected error occured while attempting to delete the user"), http.StatusNotFound)
		return
	}

	err = fmh.DB.DeleteLoansByUserID(id)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error when deleting loans by userId", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("an unexpected error occured while attempting to delete the user"), http.StatusNotFound)
		return
	}

	err = fmh.DB.DeleteCreditCardsByUserID(id)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error when deleting credit cards by userId", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("an unexpected error occured while attempting to delete the user"), http.StatusNotFound)
		return
	}

	err = fmh.DB.DeleteUserByID(id)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error when deleting user by id", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("an unexpected error occured while attempting to delete the user"), http.StatusNotFound)
		return
	}

	fmlogger.Exit(method)
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
	fmlogger.Enter(method)

	search := r.URL.Query().Get("search")
	users, err := fmh.DB.GetAllUsers(search)

	if err != nil {
		fmt.Printf("[%s] caught unexpected error when attempting to fetch users from database\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching users"), http.StatusInternalServerError)
		return
	}

	fmt.Printf("[EXIT %s]", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, users)
}
