package handlers

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (fmh *FinanceManagerHandler) CalculateLoan(w http.ResponseWriter, r *http.Request) {
	method := "loan_handler.CalculateLoan"
	fmt.Printf("[ENTER %s]\n", method)

	var payload models.Loan
	idStr := chi.URLParam(r, "userId")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("[%v] the supplied id was invalid and returned the error: %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, errors.New("invalid user id"), http.StatusBadRequest)
		return
	}

	payload.ID = id

	// Read in loan from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		fmt.Printf("[%s] failed to parse JSON payload\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	err = payload.PerformCalc()
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmt.Printf("[%s] %s\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, payload)
}

func (fmh *FinanceManagerHandler) CompareLoanPayments(w http.ResponseWriter, r *http.Request) {
	method := "loan_handler.CompareLoanPayments"
	fmt.Printf("[ENTER %s]\n", method)

	//Read ID from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
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

	// Validate that this loan belongs to the given user
	err = fmh.Validator.LoanBelongsToUser(loan, userId)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		fmt.Printf("[%s] %v", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	// Read in loan from payload
	var payload models.Loan
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		fmt.Printf("[%s] failed to read JSON payload: %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	err = loan.PerformCalc()
	err1 = payload.PerformCalc()

	if err != nil || err1 != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured during calculation"), http.StatusInternalServerError)

		if err != nil {
			fmt.Printf("[%s] error occured during calculation: %v\n", method, err)
		}

		if err1 != nil {
			fmt.Printf("[%s] error occured during calculation: %v\n", method, err1)
		}
	}

	cs := loan.CompareLoanPayments(payload)

	fmt.Printf("[EXIT %s]", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, cs)
}

func (fmh *FinanceManagerHandler) GetAllUserLoans(w http.ResponseWriter, r *http.Request) {
	method := "loan_handler.GetAllUserLoans"
	fmt.Printf("[ENTER %s]\n", method)

	//Read ID from url
	search := r.URL.Query().Get("search")
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

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

func (fmh *FinanceManagerHandler) GetLoanById(w http.ResponseWriter, r *http.Request) {
	method := "loan_handler.GetLoanById"
	fmt.Printf("[ENTER %s]\n", method)

	//Read ID from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
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

func (fmh *FinanceManagerHandler) DeleteLoanById(w http.ResponseWriter, r *http.Request) {
	method := "loan_handler.DeleteLoanById"
	fmt.Printf("[ENTER %s]\n", method)

	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
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

	// Validate that the loan exists
	loan, err := fmh.DB.GetLoanByID(loanId)
	if err != nil || loan.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmt.Printf("[%s] failed to retrieve loan", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	// Validate that the loan belongs to the user
	err = fmh.Validator.LoanBelongsToUser(loan, userId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmt.Printf("[%s] %v", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	// Delete the loan
	err = fmh.DB.DeleteLoanByID(loanId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		fmt.Printf("[%s] %v", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, "Loan deleted successfully")
}

func (fmh *FinanceManagerHandler) SaveLoan(w http.ResponseWriter, r *http.Request) {
	method := "loan_handler.SaveLoan"
	fmt.Printf("[ENTER %s]\n", method)

	var payload models.Loan

	//Read ID from url
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching loans: %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Read in loan from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		fmt.Printf("[%s] failed to read JSON payload: %v", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	payload.UserID = id

	err = payload.ValidateCanSaveLoan()
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		fmt.Printf("[%s] request is invalid: %v", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	_, err = fmh.DB.InsertLoan(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		fmt.Printf("[%s] failed to insert loan", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, "new loan was saved successfully")
}

func (fmh *FinanceManagerHandler) UpdateLoan(w http.ResponseWriter, r *http.Request) {
	method := "loan_handler.UpdateLoan"
	fmt.Printf("[ENTER %s]\n", method)

	var payload models.Loan
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
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

	// Read in loan from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		fmt.Printf("[%s] failed to parse JSON payload\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	// Validate that the loan exists
	loan, err := fmh.DB.GetLoanByID(loanId)
	if err != nil || loan.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmt.Printf("[%s] failed to retrieve loan", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	// Validate that the loan belongs to the user
	err = fmh.Validator.LoanBelongsToUser(loan, userId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmt.Printf("[%s] %v", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	// Update the loan
	err = fmh.DB.UpdateLoan(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		fmt.Printf("[%s] %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, "Loan updated successfully")
}

func (fmh *FinanceManagerHandler) GetLoanSummary(w http.ResponseWriter, r *http.Request) {
	method := "loan_handler.GetLoanSummary"
	fmt.Printf("[ENTER %s]\n", method)

	//Read ID from url
	search := r.URL.Query().Get("search")
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

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

	monthlyCost := float64(0)
	totalBalance := float64(0)
	var summary models.LoansSummary

	summary.Count = len(loans)

	for _, l := range loans {
		monthlyCost += l.MonthlyPayment
		totalBalance += l.Total
	}

	summary.MonthlyCost = monthlyCost
	summary.TotalBalance = totalBalance

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, summary)
}
