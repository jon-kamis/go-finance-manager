package fmhandler

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

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
	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, "new loan was saved successfully")
}

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
