package fmhandler

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/constants"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (fmh *FinanceManagerHandler) loadStock(ticker string) error {
	method := "stocks_handler.fetchStock"
	fmlogger.Enter(method)

	s, err := fmh.DB.GetStockByTicker(ticker)

	if err != nil {
		fmlogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	if s.ID != 0 {
		//Stock is loaded, no action needed
		fmlogger.Info(method, "stock is already loaded")
		fmlogger.Exit(method)
		return nil
	}

	s, err = fmh.StocksService.FetchStockWithTicker(ticker)

	if err != nil {
		fmlogger.ExitError(method, constants.UnexpectedExternalCallError, err)
		return err
	}

	//Persist the new stock object
	_, err = fmh.DB.InsertStock(s)

	if err != nil {
		fmlogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}

func (fmh *FinanceManagerHandler) GetIsStocksEnabled(w http.ResponseWriter, r *http.Request) {
	method := "stocks_handler.GetIsStocksEnabled"
	fmlogger.Enter(method)

	re := models.StocksEnabledResponse{
		Enabled: fmh.StocksService.GetIsStocksEnabled(),
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

	err = fmh.StocksService.UpdateAndPersistAPIKey(payload.Key)
	if err != nil {
		fmlogger.ExitError(method, constants.GenericServerError, err)
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.GenericServerError), http.StatusInternalServerError)
		return
	}

	fmh.JSONUtil.WriteJSON(w, http.StatusOK, constants.SuccessMessage)
	fmlogger.Exit(method)
}

func (fmh *FinanceManagerHandler) SaveUserStock(w http.ResponseWriter, r *http.Request) {
	method := "stocks_handler.PostStocksAPIKey"
	fmlogger.Enter(method)

	var payload models.UserStock

	//Read ID from url
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmlogger.ExitError(method, "user is not authorized to access other user data", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		return
	}

	// Read in user stock from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		fmlogger.ExitError(method, "failed to parse json payload", err)
		return
	}

	payload.UserId = id

	err = payload.ValidateCanSaveUserStock()
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		fmlogger.ExitError(method, "request is not valid", err)
		return
	}

	//Fetch or Load the requested stock
	err = fmh.loadStock(payload.Ticker)

	if err != nil {
		rerr := errors.New(constants.GenericServerError)
		fmh.JSONUtil.ErrorJSON(w, rerr, http.StatusInternalServerError)
		fmlogger.ExitError(method, constants.UnexpectedExternalCallError, err)
		return
	}

	//Insert the user stock
	_, err = fmh.DB.InsertUserStock(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		fmlogger.ExitError(method, "unexpected error occured when inserting credit card", err)
		return
	}

	fmlogger.Exit(method)
}
