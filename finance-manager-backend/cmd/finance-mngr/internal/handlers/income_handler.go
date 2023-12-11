package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (fmh *FinanceManagerHandler) GetAllUserIncomes(w http.ResponseWriter, r *http.Request) {
	method := "income_handler.GetAllUserIncomes"
	fmt.Printf("[ENTER %s]\n", method)

	//Read ID from url
	search := r.URL.Query().Get("search")
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), false, w, r)

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching loans: %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	loans, err := fmh.DB.GetAllUserLoans(id, search)

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching loans: %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching loans"), http.StatusInternalServerError)
		return
	}

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, loans)
}

func (fmh *FinanceManagerHandler) GetIncomeById(w http.ResponseWriter, r *http.Request) {
	method := "loan_handler.GetLoanById"
	fmt.Printf("[ENTER %s]\n", method)

	//Read ID from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), false, w, r)
	loanId, err1 := strconv.Atoi(chi.URLParam(r, "loanId"))

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching loans: %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err1 != nil {
		fmt.Printf("[%s] unexpected error occured when fetching loans: %v\n", method, err1)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, err1, http.StatusInternalServerError)
		return
	}

	loan, err := fmh.DB.GetLoanByID(loanId)
	if err != nil || loan.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmt.Printf("[%s] failed to retrieve loan", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	err = fmh.Validator.LoanBelongsToUser(loan, userId)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		fmt.Printf("[%s] %v", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, loan)
}
