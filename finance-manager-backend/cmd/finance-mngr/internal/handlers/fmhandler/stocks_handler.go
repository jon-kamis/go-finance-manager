package fmhandler

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/constants"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//Loads a stock from the remote API
func (fmh *FinanceManagerHandler) loadStock(ticker string) error {
	method := "stocks_handler.loadStock"
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

// SaveUserStock godoc
// @title		Insert Stock
// @version 	1.0.0
// @Tags 		Stocks
// @Summary 	Insert Stock
// @Description Inserts a new Stock into the Database for a given user
// @Param		userId path int true "User ID"
// @Param		stock body models.UserStock true "The stock to insert"
// @Accept		json
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/stocks [post]
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