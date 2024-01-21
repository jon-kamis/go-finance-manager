package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// Loads a stock from the remote API
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

	sl, err := fmh.StocksService.FetchStockWithTickerForPastYear(ticker)

	if err != nil {
		fmlogger.ExitError(method, constants.UnexpectedExternalCallError, err)
		return err
	}

	_, err = fmh.DB.InsertStock(sl[len(sl)-1])

	if err != nil {
		fmlogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	err = fmh.DB.InsertStockData(sl)

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
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, constants.SuccessMessage)
}

// GetStockHistory godoc
// @title		Get Stock History
// @version 	1.0.0
// @Tags 		Stocks
// @Summary 	Get Stock History
// @Description Gets History data for one or more stocks
// @Param		tickers query string true "A comma separated list of stocks to fetch positions for"
// @Param		days query int false "The amount of days back we should load history for. Max is 365. Default is 30"
// @Accept		json
// @Produce 	json
// @Success 	200 {array} models.PositionHistory
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/stocks [get]
func (fmh *FinanceManagerHandler) GetStockHistory(w http.ResponseWriter, r *http.Request) {
	method := "stocks_handler.GetStockHistory"
	fmlogger.Enter(method)

	tickers := r.URL.Query().Get("tickers")
	dStr := r.URL.Query().Get("days")
	var d int
	var err error
	var rArr []models.PositionHistory

	if tickers == "" {
		err = errors.New("tickers param is required")
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		fmlogger.ExitError(method, err.Error(), err)
		return
	}

	if dStr == "" {
		d = 30
	} else {
		d, err = strconv.Atoi(dStr)

		if err != nil {
			fmh.JSONUtil.ErrorJSON(w, errors.New("days param is invalid"), http.StatusBadRequest)
			fmlogger.ExitError(method, err.Error(), err)
			return
		}
	}

	historyStartDt := time.Now().Add(-1 * 24 * time.Duration(d) * time.Hour)
	tArr := strings.Split(tickers, ",")

	for _, t := range tArr {

		sd, err := fmh.DB.GetStockDataByTickerAndDateRange(t, historyStartDt, time.Now())

		if err != nil {
			fmh.JSONUtil.ErrorJSON(w, errors.New(constants.UnexpectedSQLError), http.StatusInternalServerError)
			fmlogger.ExitError(method, constants.UnexpectedSQLError, err)
			return
		}

		ph := models.PositionHistory{
			Ticker:  t,
			StartDt: historyStartDt,
			EndDt:   time.Now(),
			Count:   len(sd),
			Values:  sd,
		}

		rArr = append(rArr, ph)
	}

	fmh.JSONUtil.WriteJSON(w, http.StatusOK, rArr)
	fmlogger.Exit(method)
}
