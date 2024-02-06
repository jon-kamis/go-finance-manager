package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jon-kamis/klogger"
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
	klogger.Enter(method)

	//Read ID from url
	search := r.URL.Query().Get("search")
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		klogger.ExitError(method, "unexpected error occured when fetching bills: %v", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	bills, err := fmh.DB.GetAllUserBills(id, search)

	if err != nil {
		klogger.ExitError(method, "unexpected error occured when fetching bills: %v", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching bills"), http.StatusInternalServerError)
		return
	}

	klogger.Exit(method)
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
	klogger.Enter(method)

	//Read ID from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	billId, err1 := strconv.Atoi(chi.URLParam(r, "billId"))

	if err != nil {
		klogger.ExitError(method, "unexpected error occured when processing user id: %v", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err1 != nil {
		klogger.ExitError(method, "error occured when processing bill id: %v", err1)
		fmh.JSONUtil.ErrorJSON(w, err1, http.StatusInternalServerError)
		return
	}

	bill, err := fmh.DB.GetBillByID(billId)
	if err != nil || bill.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		klogger.ExitError(method, "error occured when fetching bill: %v", err)
		return
	}

	err = fmh.Validator.BillBelongsToUser(bill, userId)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		klogger.ExitError(method, "bill does not belong to user: %v", err)
		return
	}

	klogger.Exit(method)
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
	klogger.Enter(method)

	var payload models.Bill

	//Read ID from url
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		klogger.ExitError(method, "unexpected error occured when processing user id: %v", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Read in loan from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, "failed to parse JSON object: %v", err)
		return
	}

	payload.UserID = id

	//Validat can save bill
	err = payload.ValidateCanSaveBill()
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, "request is invalid: %v", err)
		return
	}

	_, err = fmh.DB.InsertBill(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, "error when saving income: %v", err)
		return
	}

	klogger.Exit(method)
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
	klogger.Enter(method)

	var payload models.Bill
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	billId, err1 := strconv.Atoi(chi.URLParam(r, "billId"))

	if err != nil {
		klogger.ExitError(method, "unexpected error occured when processing user id: %v", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err1 != nil {
		klogger.ExitError(method, "unexpected error occured when processing bill id: %v", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Read in bill from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, "failed to parse JSON object: %v", err)
		return
	}

	// Validate that the bill exists
	bill, err := fmh.DB.GetBillByID(billId)
	if err != nil || bill.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		klogger.ExitError(method, "bill does not exist: %v", err)
		return
	}

	// Validate that the bill belongs to the user
	err = fmh.Validator.BillBelongsToUser(bill, userId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		klogger.ExitError(method, "bill does not belong to logged in user: %v", err)
		return
	}

	//Validate the Bill object
	payload.ValidateCanSaveBill()

	// Update the bill
	err = fmh.DB.UpdateBill(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, "failed to update bill: %v", err)
		return
	}

	klogger.Exit(method)
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
	klogger.Enter(method)

	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	billId, err1 := strconv.Atoi(chi.URLParam(r, "billId"))

	if err != nil {
		klogger.ExitError(method, "unexpected error occured when processing user id: %v", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err1 != nil {
		klogger.ExitError(method, "unexpected error occured when processing bill id: %v", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Validate that the bill exists
	bill, err := fmh.DB.GetBillByID(billId)
	if err != nil || bill.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		klogger.ExitError(method, "bill does not exist", err)
		return
	}

	// Validate that the loan belongs to the user
	err = fmh.Validator.BillBelongsToUser(bill, userId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		klogger.ExitError(method, "bill does not belong to logged in user: %v", err)
		return
	}

	// Delete the loan
	err = fmh.DB.DeleteBillByID(billId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, "failed to delete bill: %v", err)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, "bill deleted successfully")
}
