package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"math"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jon-kamis/klogger"
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
	klogger.Enter(method)

	//Read ID from url
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)

		return
	}

	summary := models.Summary{}

	loans, err := fmh.DB.GetAllUserLoans(id, "")
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, "failed to retrieve user loans:\n%v", err)
		return
	}

	incomes, err := fmh.DB.GetAllUserIncomes(id, "")
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, "failed to retrieve user incomes:\n%v", err)
		return
	}

	bills, err := fmh.DB.GetAllUserBills(id, "")
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, "failed to retrieve user bills:\n%v", err)
		return
	}

	ccs, err := fmh.DB.GetAllUserCreditCards(id, "")
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, "failed to retrieve user credit cards:\n%v", err)
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

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, summary)
}

// GetUserStockPortfolioSummary godoc
// @title		Get Stock Portfolio Summary
// @version 	2.0.0
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
// @Router 		/users/{userId}/stock-portfolio [get]
func (fmh *FinanceManagerHandler) GetUserStockPortfolioSummary(w http.ResponseWriter, r *http.Request) {
	method := "summary_handler.GetUserStockPortfolioSummary"
	klogger.Enter(method)

	//Read ID from url
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	usl, err := fmh.DB.GetAllUserStocks(id, constants.UserStockTypeOwn, "", time.Now())

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.UnexpectedSQLError), http.StatusInternalServerError)
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return
	}

	//Generate list of user positions
	var pl []models.PortfolioPosition
	var sum models.UserStockPortfolioSummary

	for _, us := range usl {
		s, err := fmh.DB.GetStockByTicker(us.Ticker)

		if err != nil {
			fmh.JSONUtil.ErrorJSON(w, errors.New(constants.UnexpectedSQLError), http.StatusInternalServerError)
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return
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
		}

		pl = append(pl, p)
	}

	//Load positions into summary object
	sum.LoadPositions(pl)

	fmh.JSONUtil.WriteJSON(w, http.StatusOK, sum)
	klogger.Exit(method)
}
