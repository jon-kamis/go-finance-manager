package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/internal/finance-mngr/models/restmodels"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jon-kamis/klogger"
)

// Loads a stock from the remote API
func (fmh *FinanceManagerHandler) loadStock(ticker string) error {
	method := "stocks_handler.loadStock"
	klogger.Enter(method)

	s, err := fmh.DB.GetStockByTicker(ticker)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	if s.ID != 0 {
		//Stock is loaded, no action needed
		klogger.Info(method, "stock is already loaded")
		klogger.Exit(method)
		return nil
	}

	sl, err := fmh.ExternalService.FetchStockWithTickerForPastYear(ticker)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedExternalCallError, err)
		return err
	}

	_, err = fmh.DB.InsertStock(sl[len(sl)-1])

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	err = fmh.DB.InsertStockData(sl)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	klogger.Exit(method)
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
	klogger.Enter(method)

	var payload models.UserStock

	//Read ID from url
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		return
	}

	// Read in user stock from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, constants.FailedToParseJsonBodyError, err)
		return
	}

	payload.UserId = id

	err = payload.ValidateCanSaveUserStock()
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, constants.GenericBadRequestErrorLog, err)
		return
	}

	//Fetch or Load the requested stock
	err = fmh.loadStock(payload.Ticker)

	if err != nil {
		rerr := errors.New(constants.GenericServerError)
		fmh.JSONUtil.ErrorJSON(w, rerr, http.StatusInternalServerError)
		klogger.ExitError(method, constants.UnexpectedExternalCallError, err)
		return
	}

	//Insert the user stock
	_, err = fmh.DB.InsertUserStock(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.FailedToSaveEntityError, err)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, constants.SuccessMessage)
}

// ModifyStockOperation godoc
// @title		Modify Stock Operation
// @version 	1.0.0
// @Tags 		Stocks
// @Summary 	Modify User Stock
// @Description Modifies a user's stock. This is an add or remove operation and can be used to post new stock
// @Param		userId path int true "ID of the user to modify stocks for"
// @Param		request body restmodels.ModifyStockRequest true "The request to process"
// @Accept		json
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/stock-operation [post]
func (fmh *FinanceManagerHandler) ModifyStockOperation(w http.ResponseWriter, r *http.Request) {
	method := "stocks_handler.ModifyStockOperation"
	klogger.Enter(method)

	var p restmodels.ModifyStockRequest

	//Read UserId from url
	uId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		return
	}

	// Read in request from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &p)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, constants.FailedToParseJsonBodyError, err)
		return
	}

	// Validate Request
	isValid, errMsg := p.IsValidRequest()

	if !isValid {
		err = errors.New(errMsg)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, errMsg)
		return
	}

	//Load stock if required
	err = fmh.loadStock(p.Ticker)

	if err != nil {
		rerr := errors.New(constants.GenericServerError)
		fmh.JSONUtil.ErrorJSON(w, rerr, http.StatusInternalServerError)
		klogger.ExitError(method, constants.UnexpectedExternalCallError, err)
		return
	}

	//Prior user stock
	var usp models.UserStock

	us := models.UserStock{
		UserId:      uId,
		Ticker:      p.Ticker,
		Type:        constants.UserStockTypeOwn,
		EffectiveDt: p.Date,
	}

	err = fmh.Service.LoadPriorUserStockForTransaction(p, &usp, &us)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.UnexpectedExternalCallError, err)
		return
	}

	//Update usp if it exists
	if usp.ID != 0 {
		err = fmh.DB.UpdateUserStock(usp)

		if err != nil {
			fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
			klogger.ExitError(method, constants.FailedToUpdateEntityError, err)
			return
		}
	}

	//Save new user stock created by operation if quantity is greater than 0
	if us.Quantity > 0 {
		_, err = fmh.DB.InsertUserStock(us)

		if err != nil {
			fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
			klogger.ExitError(method, constants.FailedToSaveEntityError, err)
			return
		}
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, constants.SuccessMessage)
}

// GetUserStocks godoc
// @title		Get User Stocks
// @version 	1.0.0
// @Tags 		Stocks
// @Summary 	Get User Stocks
// @Description Gets a list of stocks currently owned or watched by a given user
// @Param		userId path string true "The ID of the user to fetch stocks for"
// @Param		type query string false "The Type of stocks to fetch. Available types are {'own','watchlist'}. Default is 'own'."
// @Param		search query string false "Search for User stocks by ticker"
// @Accept		json
// @Produce 	json
// @Success 	200 {array} models.Stock
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/stocks [get]
func (fmh *FinanceManagerHandler) GetUserStocks(w http.ResponseWriter, r *http.Request) {
	method := "stocks_handler.GetUserStocks"
	klogger.Enter(method)

	//Read ID from url
	uid, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	stockType := r.URL.Query().Get("type")
	search := r.URL.Query().Get("search")
	var sl []models.Stock

	usl, err := fmh.DB.GetAllUserStocks(uid, stockType, search, time.Now())

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.FailedToRetrieveEntityError, err)
		return
	}

	for _, us := range usl {
		s, err := fmh.DB.GetStockByTicker(us.Ticker)

		if err != nil {
			fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return
		}

		sl = append(sl, s)
	}

	fmh.JSONUtil.WriteJSON(w, http.StatusOK, sl)
	klogger.Exit(method)
}

// GetStockHistory godoc
// @title		Get Stock History
// @version 	2.1.0
// @Tags 		Stocks
// @Summary 	Get Stock History
// @Description Gets History data for one or more stocks
// @Param		tickers query string true "A comma separated list of stocks to fetch positions for"
// @Param		histLength query int false "The lenght of history to fetch. Available values are 'day', 'week', 'month', and 'year'. Default is 'month'"
// @Accept		json
// @Produce 	json
// @Success 	200 {array} models.PositionHistory
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/stocks [get]
func (fmh *FinanceManagerHandler) GetStockHistory(w http.ResponseWriter, r *http.Request) {
	method := "stocks_handler.GetStockHistory"
	klogger.Enter(method)

	tickers := r.URL.Query().Get("tickers")
	hlStr := r.URL.Query().Get("histLength")

	var d int
	var err error
	var rArr []models.PositionHistory

	if tickers == "" {
		err = errors.New("tickers param is required")
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, constants.GenericBadRequestErrorLog, err)
		return
	}

	if hlStr == "" {
		hlStr = constants.LengthWeek
	}

	switch hlStr {
	case constants.LengthDay:
		d = 1
	case constants.LengthWeek:
		d = 7
	case constants.LengthMonth:
		d = 31
	case constants.LengthYear:
		d = 365
	default:
		d = 31
	}

	historyStartDt := time.Now().Add(-1 * 24 * time.Duration(d) * time.Hour)
	tArr := strings.Split(tickers, ",")

	for _, t := range tArr {

		sd, err := fmh.DB.GetStockDataByTickerAndDateRange(t, historyStartDt, time.Now())

		if err != nil {
			fmh.JSONUtil.ErrorJSON(w, errors.New(constants.UnexpectedSQLError), http.StatusInternalServerError)
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return
		}

		var high float64
		var low float64
		var open float64
		var close float64
		var delta float64
		var deltaPercentage float64

		high = sd[0].High
		low = sd[0].Low
		open = sd[0].Open
		close = sd[len(sd)-1].Close
		delta = sd[len(sd)-1].Close - sd[0].Open
		deltaPercentage = delta / sd[0].Open * 100

		//Populate high and low
		for _, s := range sd {
			if s.High > high {
				high = s.High
			}

			if s.Low < low {
				low = s.Low
			}
		}

		ph := models.PositionHistory{
			Ticker:          t,
			High:            high,
			Low:             low,
			Open:            open,
			Close:           close,
			Delta:           delta,
			DeltaPercentage: deltaPercentage,
			StartDt:         historyStartDt,
			EndDt:           time.Now(),
			Count:           len(sd),
			Values:          sd,
		}

		rArr = append(rArr, ph)
	}

	fmh.JSONUtil.WriteJSON(w, http.StatusOK, rArr)
	klogger.Exit(method)
}

// GetStockHistory godoc
// @title		Get User Stock Portfolio History
// @version 	1.0.0
// @Tags 		Stocks
// @Summary 	Get User Stock Portfolio History
// @Description Gets History of a User's Stock Portfolio Balance
// @Param		userId path int true "The ID of the user to get Portfolio History for"
// @Param		histLength query int false "The lenght of history to fetch. Available values are 'week', 'month', and 'year'. Default is 'week'"
// @Accept		json
// @Produce 	json
// @Success 	200 {object} models.StockPortfolioHistoryResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/stock-portfolio-history [get]
func (fmh *FinanceManagerHandler) GetUserStockPortfolioHistory(w http.ResponseWriter, r *http.Request) {
	method := "stocks_handler.GetUserStockPortfolioHistory"
	klogger.Enter(method)

	var resp models.StockPortfolioHistoryResponse
	var err error
	var hl int

	//Read URL variables
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	hlStr := r.URL.Query().Get("histLength")

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	if hlStr == "" {
		hlStr = constants.LengthWeek
	}

	switch hlStr {
	case constants.LengthWeek:
		hl = 7
	case constants.LengthMonth:
		hl = 31
	case constants.LengthYear:
		hl = 365
	default:
		hl = 7
	}

	//Load positions History object
	h, err := fmh.Service.GetUserPortfolioBalanceHistory(id, hl)
	resp.Items = h
	resp.Count = len(h)

	//Get highest and lowest value
	high := h[0].High
	low := h[0].Low

	for _, i := range h {
		if i.Low < low {
			low = i.Low
		}

		if i.High > high {
			high = i.High
		}
	}

	resp.High = high
	resp.Low = low
	resp.Open = h[0].Open
	resp.Close = h[len(h)-1].Close
	resp.Delta = resp.Close - resp.Open
	resp.DeltaPercentage = resp.Delta / resp.Open * 100

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.GenericServerError), http.StatusInternalServerError)
		klogger.ExitError(method, constants.GenericServerError, err)
		return
	}

	fmh.JSONUtil.WriteJSON(w, http.StatusOK, resp)
	klogger.Exit(method)
}
