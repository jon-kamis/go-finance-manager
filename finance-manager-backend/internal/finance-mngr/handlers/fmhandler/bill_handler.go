package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetAllUserBills godoc
// @title		 Get All User Bills
// @version 	1.0.0
// @Tags 		Bills
// @Summary 	Get All User Bills
// @Description Returns an array of Bill objects belonging to a given user
// @Param		userId path int true "User ID"
// @Param		search query string false "Search for bills by name"
// @Produce 	json
// @Success 	200 {array} models.Bill
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/bills [get]
func (fmh *FinanceManagerHandler) GetAllUserBills(w http.ResponseWriter, r *http.Request) {
	method := "bill_handler.GetAllUserBills"
	fmlogger.Enter(method)

	//Read ID from url
	search := r.URL.Query().Get("search")
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching bills", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	bills, err := fmh.DB.GetAllUserBills(id, search)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching bills", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching incomes"), http.StatusInternalServerError)
		return
	}

	fmlogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, bills)
}

// GetBillById godoc
// @title		Get User Bill by ID
// @version 	1.0.0
// @Tags 		Bills
// @Summary 	Get Bill by ID
// @Description Returns a Bill by its ID for a given user
// @Param		userId path int true "User ID"
// @Param		billId path int true "Bill ID"
// @Produce 	json
// @Success 	200 {object} models.Bill
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/bills/{billId} [get]
func (fmh *FinanceManagerHandler) GetBillById(w http.ResponseWriter, r *http.Request) {
	method := "bill_handler.GetBillById"
	fmlogger.Enter(method)

	//Read ID from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	billId, err1 := strconv.Atoi(chi.URLParam(r, "billId"))

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when processing user id", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err1 != nil {
		fmlogger.ExitError(method, "error occured when processing bill id", err1)
		fmh.JSONUtil.ErrorJSON(w, err1, http.StatusInternalServerError)
		return
	}

	bill, err := fmh.DB.GetBillByID(billId)
	if err != nil || bill.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmlogger.ExitError(method, "error occured when fetching bill", err)
		return
	}

	err = fmh.Validator.BillBelongsToUser(bill, userId)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		fmlogger.ExitError(method, "bill does not belong to user", err)
		return
	}

	fmlogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, bill)
}

// SaveBill godoc
// @title		Insert Bill
// @version 	1.0.0
// @Tags 		Bills
// @Summary 	Insert Bill
// @Description Inserts a new Bill into the Database for a given user
// @Param		userId path int true "User ID"
// @Param		bill body models.Bill true "The bill to insert"
// @Accept		json
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/bills [post]
func (fmh *FinanceManagerHandler) SaveBill(w http.ResponseWriter, r *http.Request) {
	method := "bill_handler.SaveBill"
	fmlogger.Enter(method)

	var payload models.Bill

	//Read ID from url
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when processing user id", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Read in loan from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		fmlogger.ExitError(method, "failed to parse JSON object", err)
		return
	}

	payload.UserID = id

	//Validat can save bill
	err = payload.ValidateCanSaveBill()
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		fmlogger.ExitError(method, "request is invalid", err)
		return
	}

	_, err = fmh.DB.InsertBill(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		fmlogger.ExitError(method, "error when saving income", err)
		return
	}

	fmlogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, constants.SuccessMessage)
}

// UpdateBill godoc
// @title		Update Bill
// @version 	1.0.0
// @Tags 		Bills
// @Summary 	Update Bill
// @Description Updates an existing Bill for a user
// @Param		userId path int true "User ID"
// @Param		billId path int true "ID of the bill to update"
// @Param		bill body models.Bill true "The bill to update"
// @Accept		json
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/bills/{billId} [put]
func (fmh *FinanceManagerHandler) UpdateBill(w http.ResponseWriter, r *http.Request) {
	method := "bill_handler.UpdateBill"
	fmlogger.Enter(method)

	var payload models.Bill
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	billId, err1 := strconv.Atoi(chi.URLParam(r, "billId"))

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when processing user id", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err1 != nil {
		fmlogger.ExitError(method, "unexpected error occured when processing bill id", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Read in bill from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		fmlogger.ExitError(method, "failed to parse JSON object", err)
		return
	}

	// Validate that the bill exists
	bill, err := fmh.DB.GetBillByID(billId)
	if err != nil || bill.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmlogger.ExitError(method, "bill does not exist", err)
		return
	}

	// Validate that the bill belongs to the user
	err = fmh.Validator.BillBelongsToUser(bill, userId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		fmlogger.ExitError(method, "bill does not belong to logged in user", err)
		return
	}

	//Validate the Bill object
	payload.ValidateCanSaveBill()

	// Update the bill
	err = fmh.DB.UpdateBill(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		fmlogger.ExitError(method, "failed to update bill", err)
		return
	}

	fmlogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, "Bill updated successfully")
}

// DeleteBill godoc
// @title		Delete Bill
// @version 	1.0.0
// @Tags 		Bills
// @Summary 	Delete Bill by ID
// @Description Deletes a user's Bill by its ID
// @Param		userId path int true "User ID"
// @Param		billId path int true "ID of the bill"
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/bills/{billId} [delete]
func (fmh *FinanceManagerHandler) DeleteBillById(w http.ResponseWriter, r *http.Request) {
	method := "bill_handler.DeleteBillById"
	fmlogger.Enter(method)

	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	billId, err1 := strconv.Atoi(chi.URLParam(r, "billId"))

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when processing user id", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err1 != nil {
		fmlogger.ExitError(method, "unexpected error occured when processing bill id", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Validate that the bill exists
	bill, err := fmh.DB.GetBillByID(billId)
	if err != nil || bill.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		fmlogger.ExitError(method, "bill does not exist", err)
		return
	}

	// Validate that the loan belongs to the user
	err = fmh.Validator.BillBelongsToUser(bill, userId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		fmlogger.ExitError(method, "bill does not belong to logged in user", err)
		return
	}

	// Delete the loan
	err = fmh.DB.DeleteBillByID(billId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		fmlogger.ExitError(method, "failed to delete bill", err)
		return
	}

	fmlogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, "bill deleted successfully")
}
