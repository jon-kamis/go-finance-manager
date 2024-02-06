package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jon-kamis/klogger"
)

// GetAllUserIncomes godoc
// @title		Get All User Incomes
// @version 	1.0.0
// @Tags 		Incomes
// @Summary 	Get All User Incomes
// @Description Returns an array of Income objects belonging to a given user
// @Param		userId path int true "User ID"
// @Param		search query string false "Search for incomes by name"
// @Produce 	json
// @Success 	200 {array} models.Income
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/incomes [get]
func (fmh *FinanceManagerHandler) GetAllUserIncomes(w http.ResponseWriter, r *http.Request) {
	method := "income_handler.GetAllUserIncomes"
	klogger.Enter(method)

	//Read ID from url
	search := r.URL.Query().Get("search")
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	incomes, err := fmh.DB.GetAllUserIncomes(id, search)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching incomes"), http.StatusInternalServerError)
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return
	}

	for _, i := range incomes {
		i.PopulateEmptyValues(time.Now())
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, incomes)
}

// GetIncomeById godoc
// @title		Get Income by ID
// @version 	1.0.0
// @Tags 		Incomes
// @Summary 	Get Income by ID
// @Description Returns an Income object belonging to a given user
// @Param		userId path int true "User ID"
// @Param		incomeId path int true "the ID of the Income"
// @Produce 	json
// @Success 	200 {object} models.Income
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/incomes/{incomeId} [get]
func (fmh *FinanceManagerHandler) GetIncomeById(w http.ResponseWriter, r *http.Request) {
	method := "income_handler.GetIncomeById"
	klogger.Enter(method)

	//Read ID from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	incomeId, err1 := strconv.Atoi(chi.URLParam(r, "incomeId"))

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	if err1 != nil {
		fmh.JSONUtil.ErrorJSON(w, err1, http.StatusInternalServerError)
		klogger.ExitError(method, constants.ProcessIdError, err)
		return
	}

	income, err := fmh.DB.GetIncomeByID(incomeId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.UnexpectedSQLError), http.StatusInternalServerError)
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return
	}

	if income.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusNotFound)
		klogger.ExitError(method, constants.EntityNotFoundError)
		return
	}

	err = fmh.Validator.IncomeBelongsToUser(income, userId)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	err = income.PopulateEmptyValues(time.Now())
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		klogger.ExitError(method, constants.GenericUnprocessableEntityErrLog, err)
		return
	}

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, income)
}

// SaveIncome godoc
// @title		Insert Income
// @version 	1.0.0
// @Tags 		Incomes
// @Summary 	Insert Income
// @Description Inserts a new Income into the Database for a given user
// @Param		userId path int true "User ID"
// @Param		income body models.Income true "The income to insert"
// @Accept		json
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/incomes [post]
func (fmh *FinanceManagerHandler) SaveIncome(w http.ResponseWriter, r *http.Request) {
	method := "income_handler.SaveIncome"
	klogger.Enter(method)

	var payload models.Income

	//Read ID from url
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	// Read in loan from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, constants.FailedToParseJsonBodyError, err)
		return
	}

	payload.UserID = id

	payload.PopulateEmptyValues(time.Now())

	err = payload.ValidateCanSaveIncome()
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, constants.GenericUnprocessableEntityErrLog, err)
		return
	}

	_, err = fmh.DB.InsertIncome(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.FailedToSaveEntityError, err)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, "new loan was saved successfully")
}

// UpdateIncome godoc
// @title		Update Income
// @version 	1.0.0
// @Tags 		Incomes
// @Summary 	Update Income
// @Description Updates an existing Income for a user
// @Param		userId path int true "User ID"
// @Param		incomeId path int true "ID of the income to update"
// @Param		income body models.Income true "The income to update"
// @Accept		json
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/incomes/{incomeId} [put]
func (fmh *FinanceManagerHandler) UpdateIncome(w http.ResponseWriter, r *http.Request) {
	method := "income_handler.UpdateIncome"
	klogger.Enter(method)

	var payload models.Income
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	incomeId, err1 := strconv.Atoi(chi.URLParam(r, "incomeId"))

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	if err1 != nil {
		fmh.JSONUtil.ErrorJSON(w, err1, http.StatusInternalServerError)
		klogger.ExitError(method, constants.ProcessIdError, err)
		return
	}

	// Read in income from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, constants.FailedToParseJsonBodyError, err)
		return
	}

	// Validate that the income exists
	income, err := fmh.DB.GetIncomeByID(incomeId)
	if err != nil || income.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusNotFound)
		klogger.ExitError(method, constants.EntityNotFoundError)
		return
	}

	// Validate that the income belongs to the user
	err = fmh.Validator.IncomeBelongsToUser(income, userId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
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
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, "Loan updated successfully")
}

// DeleteIncomeById godoc
// @title		Delete Income
// @version 	1.0.0
// @Tags 		Incomes
// @Summary 	Delete Income by ID
// @Description Deletes a user's Income by its ID
// @Param		userId path int true "User ID"
// @Param		incomeId path int true "ID of the Income"
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/incomes/{incomeId} [delete]
func (fmh *FinanceManagerHandler) DeleteIncomeById(w http.ResponseWriter, r *http.Request) {
	method := "income_handler.DeleteIncomeById"
	klogger.Enter(method)

	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	incomeId, err1 := strconv.Atoi(chi.URLParam(r, "incomeId"))

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.BillDoesNotBelongToUserError, err)
		return
	}

	if err1 != nil {
		fmh.JSONUtil.ErrorJSON(w, err1, http.StatusInternalServerError)
		klogger.ExitError(method, constants.ProcessIdError, err)
		return
	}

	// Validate that the loan exists
	income, err := fmh.DB.GetIncomeByID(incomeId)
	if err != nil || income.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusNotFound)
		klogger.ExitError(method, constants.EntityNotFoundError)
		return
	}

	// Validate that the loan belongs to the user
	err = fmh.Validator.IncomeBelongsToUser(income, userId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError)
		return
	}

	// Delete the loan
	err = fmh.DB.DeleteIncomeByID(incomeId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.FailedToDeleteEntityError, err)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, "Loan deleted successfully")
}
