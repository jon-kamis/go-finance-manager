package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// GetUserSummary godoc
// @title		Get Finance Summary
// @version 	1.0.0
// @Tags 		Summary
// @Summary 	Get Finance Summary
// @Description Gets a summary of all financial data for a user
// @Param		userId path int true "User ID"
// @Accept		json
// @Produce 	json
// @Success 	200 {object} models.Summary
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/summary [get]
func (fmh *FinanceManagerHandler) GetUserSummary(w http.ResponseWriter, r *http.Request) {
	method := "summary_handler.GetUserSummary"
	fmlogger.Enter(method)

	//Read ID from url
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when reading url parameters", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	summary := models.Summary{}

	loans, err := fmh.DB.GetAllUserLoans(id, "")
	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching loans", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	incomes, err := fmh.DB.GetAllUserIncomes(id, "")
	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching incomes", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	bills, err := fmh.DB.GetAllUserBills(id, "")
	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching bills", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	ccs, err := fmh.DB.GetAllUserCreditCards(id, "")
	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching credit cards", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	for _, i := range incomes {
		i.PopulateEmptyValues(time.Now())
	}

	summary.LoadLoans(loans)
	summary.LoadIncomes(incomes)
	summary.LoadBills(bills)
	summary.LoadCreditCards(ccs)

	summary.Finalize()

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, summary)
}

// GetUserStockPortfolioSummary godoc
// @title		Get Stock Portfolio Summary
// @version 	1.0.0
// @Tags 		Summary
// @Summary 	Get Stock Portfolio Summary
// @Description Gets a summary of all stock data for a user
// @Param		userId path int true "User ID"
// @Accept		json
// @Produce 	json
// @Success 	200 {object} models.UserStockPortfolioSummary
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/stocks [get]
func (fmh *FinanceManagerHandler) GetUserStockPortfolioSummary(w http.ResponseWriter, r *http.Request) {
	method := "summary_handler.GetUserStockPortfolioSummary"
	fmlogger.Enter(method)

	//Read ID from url
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when reading url parameters", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		return
	}

	usl, err := fmh.DB.GetAllUserStocks(id, "", time.Now())

	if err != nil {
		fmlogger.ExitError(method, constants.UnexpectedSQLError, err)
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.UnexpectedSQLError), http.StatusInternalServerError)
		return
	}

	//Generate list of user positions
	var pl []models.PortfolioPosition
	var sum models.UserStockPortfolioSummary
	historyStartDt := time.Now().Add(-1 * 24 * 30 * time.Hour)

	for _, us := range usl {
		s, err := fmh.DB.GetStockByTicker(us.Ticker)

		if err != nil {
			fmh.JSONUtil.ErrorJSON(w, errors.New(constants.UnexpectedSQLError), http.StatusInternalServerError)
			fmlogger.ExitError(method, constants.UnexpectedSQLError, err)
			return
		}

		sd, err := fmh.DB.GetStockDataByTickerAndDateRange(s.Ticker, historyStartDt, time.Now())

		if err != nil {
			fmh.JSONUtil.ErrorJSON(w, errors.New(constants.UnexpectedSQLError), http.StatusInternalServerError)
			fmlogger.ExitError(method, constants.UnexpectedSQLError, err)
			return
		}

		ph := models.PositionHistory{
			Ticker:  us.Ticker,
			StartDt: historyStartDt,
			EndDt:   time.Now(),
			Count:   len(sd),
			Values:  sd,
		}

		p := models.PortfolioPosition{
			Ticker:   us.Ticker,
			Quantity: us.Quantity,
			Value:    math.Round(s.Close*us.Quantity*100) / 100,
			Open:     s.Open,
			Close:    s.Close,
			High:     s.High,
			Low:      s.Low,
			AsOfDate: s.Date,
			History:  ph,
		}

		pl = append(pl, p)
	}

	//Load positions into summary object
	sum.LoadPositions(pl)

	fmh.JSONUtil.WriteJSON(w, http.StatusOK, sum)
	fmlogger.Exit(method)
}
