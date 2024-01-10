package fmhandler

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

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
