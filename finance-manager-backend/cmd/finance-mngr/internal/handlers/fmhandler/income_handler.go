package fmhandler

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func (fmh *FinanceManagerHandler) GetAllUserIncomes(w http.ResponseWriter, r *http.Request) {
	method := "income_handler.GetAllUserIncomes"
	fmt.Printf("[ENTER %s]\n", method)

	//Read ID from url
	search := r.URL.Query().Get("search")
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching incomes: %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	incomes, err := fmh.DB.GetAllUserIncomes(id, search)

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching incomes: %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching incomes"), http.StatusInternalServerError)
		return
	}

	for _, i := range incomes {
		i.PopulateEmptyValues(time.Now())
	}

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, incomes)
}

func (fmh *FinanceManagerHandler) GetIncomeById(w http.ResponseWriter, r *http.Request) {
	method := "income_handler.GetIncomeById"
	fmt.Printf("[ENTER %s]\n", method)

	//Read ID from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	incomeId, err1 := strconv.Atoi(chi.URLParam(r, "incomeId"))

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching income: %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err1 != nil {
		fmt.Printf("[%s] unexpected error occured when fetching income: %v\n", method, err1)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, err1, http.StatusInternalServerError)
		return
	}

	income, err := fmh.DB.GetIncomeByID(incomeId)
	if err != nil || income.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmt.Printf("[%s] failed to retrieve income", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	err = fmh.Validator.IncomeBelongsToUser(income, userId)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		fmt.Printf("[%s] %v", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	err = income.PopulateEmptyValues(time.Now())
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmt.Printf("[%s] %v", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, income)
}

func (fmh *FinanceManagerHandler) SaveIncome(w http.ResponseWriter, r *http.Request) {
	method := "income_handler.SaveIncome"
	fmt.Printf("[ENTER %s]\n", method)

	var payload models.Income

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

	payload.PopulateEmptyValues(time.Now())

	err = payload.ValidateCanSaveIncome()
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		fmt.Printf("[%s] request is invalid: %v", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	_, err = fmh.DB.InsertIncome(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		fmt.Printf("[%s] failed to insert income", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, "new loan was saved successfully")
}

func (fmh *FinanceManagerHandler) UpdateIncome(w http.ResponseWriter, r *http.Request) {
	method := "income_handler.UpdateIncome"
	fmt.Printf("[ENTER %s]\n", method)

	var payload models.Income
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	incomeId, err1 := strconv.Atoi(chi.URLParam(r, "incomeId"))

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching income: %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err1 != nil {
		fmt.Printf("[%s] unexpected error occured when fetching income: %v\n", method, err1)
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
	income, err := fmh.DB.GetIncomeByID(incomeId)
	if err != nil || income.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmt.Printf("[%s] failed to retrieve income", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	// Validate that the income belongs to the user
	err = fmh.Validator.IncomeBelongsToUser(income, userId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmt.Printf("[%s] %v", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	//Populate Income values
	payload.PopulateEmptyValues(time.Now())

	//Validate the Income object
	payload.ValidateCanSaveIncome()

	// Update the loan
	err = fmh.DB.UpdateIncome(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		fmt.Printf("[%s] %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, "Loan updated successfully")
}

func (fmh *FinanceManagerHandler) DeleteIncomeById(w http.ResponseWriter, r *http.Request) {
	method := "income_handler.DeleteIncomeById"
	fmt.Printf("[ENTER %s]\n", method)

	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	incomeId, err1 := strconv.Atoi(chi.URLParam(r, "incomeId"))

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching income: %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err1 != nil {
		fmt.Printf("[%s] unexpected error occured when fetching income: %v\n", method, err1)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, err1, http.StatusInternalServerError)
		return
	}

	// Validate that the loan exists
	income, err := fmh.DB.GetIncomeByID(incomeId)
	if err != nil || income.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmt.Printf("[%s] failed to retrieve income", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	// Validate that the loan belongs to the user
	err = fmh.Validator.IncomeBelongsToUser(income, userId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmt.Printf("[%s] %v", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	// Delete the loan
	err = fmh.DB.DeleteIncomeByID(incomeId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		fmt.Printf("[%s] %v", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, "Loan deleted successfully")
}
