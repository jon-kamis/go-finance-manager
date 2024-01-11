package fmhandler

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/constants"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"net/http"
)

func (fmh *FinanceManagerHandler) GetIsStocksEnabled(w http.ResponseWriter, r *http.Request) {
	method := "stocks_handler.GetIsStocksEnabled"
	fmlogger.Enter(method)

	re := models.StocksEnabledResponse{
		Enabled: fmh.StocksEnabled,
	}

	fmh.JSONUtil.WriteJSON(w, http.StatusOK, re)
	fmlogger.Exit(method)
}

func (fmh *FinanceManagerHandler) PostStocksAPIKey(w http.ResponseWriter, r *http.Request) {
	method := "stocks_handler.PostStocksAPIKey"
	fmlogger.Enter(method)

	uId, err := fmh.Auth.GetLoggedInUserId(w, r)

	//uId must be loaded successfully to proceed
	if err != nil {
		fmlogger.ExitError(method, constants.FailedToReadUserIdFromAuthHeaderError, err)
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.FailedToReadUserIdFromAuthHeaderError), http.StatusInternalServerError)
		return
	}

	//Determine if user has admin role
	hasRole, err := fmh.Validator.CheckIfUserHasRole(uId, "admin")

	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	//User must be admin to proceed
	if !hasRole {
		err = errors.New(constants.GenericForbiddenError)
		fmlogger.ExitError(method, err.Error(), err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		return
	}

	var payload models.EnableStocksRequest

	// Read payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmlogger.ExitError(method, constants.GenericBadRequestError, err)
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.GenericBadRequestError), http.StatusBadRequest)
		return
	}

	err = fmh.UpdateAndPersistAPIKey(payload.Key)
	if err != nil {
		fmlogger.ExitError(method, constants.GenericServerError, err)
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.GenericServerError), http.StatusInternalServerError)
		return
	}

	fmh.JSONUtil.WriteJSON(w, http.StatusOK, constants.SuccessMessage)
	fmlogger.Exit(method)
}
