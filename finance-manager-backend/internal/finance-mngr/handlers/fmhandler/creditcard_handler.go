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

// SaveCreditcard godoc
// @title		Insert Credit Card
// @version 	1.0.0
// @Tags 		Credit Cards
// @Summary 	Insert Credit Card
// @Description Inserts a new Credit Card object for a given user
// @Param		userId path int true "User ID"
// @Param       creditCard body models.CreditCard true "Credit card to insert"
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/credit-cards [post]
func (fmh *FinanceManagerHandler) SaveCreditCard(w http.ResponseWriter, r *http.Request) {
	method := "creditcard_handler.SaveCreditCard"
	klogger.Enter(method)

	var payload models.CreditCard

	//Read ID from url
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		klogger.ExitError(method, "user is not authorized to view other user data: %v", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		return
	}

	// Read in loan from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, "failed to parse json payload: %v", err)
		return
	}

	payload.UserID = id

	err = payload.ValidateCanSaveCreditCard()
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, "credit card request is not valid: %v", err)
		return
	}

	_, err = fmh.DB.InsertCreditCard(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, "unexpected error occured when inserting credit card: %v", err)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, "success")
}

// GetAllUserCreditCards godoc
// @title		 Get All User Credit Cards
// @version 	1.0.0
// @Tags 		Credit Cards
// @Summary 	Get All User Credit Cards
// @Description Returns an array of CreditCard objects belonging to a given user
// @Param		userId path int true "User ID"
// @Param		search query string false "Search for Credit Cards by name"
// @Produce 	json
// @Success 	200 {array} models.CreditCard
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/credit-cards [get]
func (fmh *FinanceManagerHandler) GetAllUserCreditCards(w http.ResponseWriter, r *http.Request) {
	method := "creditcard_handler.GetAllUserCreditCards"
	klogger.Enter(method)

	//Read ID from url
	search := r.URL.Query().Get("search")
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		klogger.ExitError(method, "unexpected error occured when fetching credit cards: %v", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		return
	}

	ccs, err := fmh.DB.GetAllUserCreditCards(id, search)

	//Calculate All Payments
	for _, cc := range ccs {
		cc.CalcPayment()
	}

	if err != nil {
		klogger.ExitError(method, "unexpected error occured when fetching credit cards: %v", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.UnexpectedSQLError), http.StatusInternalServerError)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, ccs)
}

// GetCreditCardById godoc
// @title		Get Credit Card by ID
// @version 	1.0.0
// @Tags 		Credit Cards
// @Summary 	Get Credit Card by ID
// @Description Fetches a Credit Card by its ID for a given user
// @Param		userId path int true "User ID"
// @Param		ccId path int true "ID of the Credit Card"
// @Produce 	json
// @Success 	200 {object} models.CreditCard
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/credit-cards/{ccId} [get]
func (fmh *FinanceManagerHandler) GetCreditCardById(w http.ResponseWriter, r *http.Request) {
	method := "creditcard_handler.GetCreditCardById"
	klogger.Enter(method)

	//Read ID from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	ccId, err1 := strconv.Atoi(chi.URLParam(r, "ccId"))

	if err != nil {
		klogger.ExitError(method, constants.ProcessUserIdError, err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		return
	}

	if err1 != nil {
		klogger.ExitError(method, constants.ProcessIdError, err1)
		fmh.JSONUtil.ErrorJSON(w, err1, http.StatusBadRequest)
		return
	}

	cc, err := fmh.DB.GetCreditCardByID(ccId)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return
	}

	if cc.ID == 0 {
		err = errors.New(constants.EntityNotFoundError)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusNotFound)
		klogger.ExitError(method, constants.EntityNotFoundError, err)
		return
	}

	err = fmh.Validator.CreditCardBelongsToUser(cc, userId)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		klogger.ExitError(method, constants.UserForbiddenToViewOtherUserDataError, err)
		return
	}

	//Calculate Payment
	cc.CalcPayment()

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, cc)
}

// DeleteCreditCardById godoc
// @title		Delete Credit Card by ID
// @version 	1.0.0
// @Tags 		Credit Cards
// @Summary 	Delete Credit Card by ID
// @Description Deletes a Credit Card by its ID for a given user
// @Param		userId path int true "User ID"
// @Param		ccId path int true "ID of the Credit Card"
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/credit-cards/{ccId} [delete]
func (fmh *FinanceManagerHandler) DeleteCreditCardById(w http.ResponseWriter, r *http.Request) {
	method := "creditcard_handler.DeleteCreditCardById"
	klogger.Enter(method)

	_, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	ccId, err1 := strconv.Atoi(chi.URLParam(r, "ccId"))

	if err != nil {
		klogger.ExitError(method, constants.ProcessUserIdError, err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		return
	}

	if err1 != nil {
		klogger.ExitError(method, constants.ProcessIdError, err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate that the bill exists
	cc, err := fmh.DB.GetCreditCardByID(ccId)
	if err != nil || cc.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusNotFound)
		klogger.ExitError(method, "credit card does not exist", err)
		return
	}

	// Get loggedIn userId
	loggedInUserId, err := fmh.Auth.GetLoggedInUserId(w, r)

	if err != nil {
		klogger.ExitError(method, constants.FailedToReadUserIdFromAuthHeaderError, err)
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.GenericServerError), http.StatusInternalServerError)
		return
	}

	// Validate that the loan belongs to the user
	err = fmh.Validator.CreditCardBelongsToUser(cc, loggedInUserId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	// Delete the loan
	err = fmh.DB.DeleteCreditCardsByID(ccId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.FailedToDeleteEntityError, err)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, constants.SuccessMessage)
}

// UpdateCreditCard godoc
// @title		Update Credit Card by ID
// @version 	1.0.0
// @Tags 		Credit Cards
// @Summary 	Update Credit Card by ID
// @Description Updates a Credit Card by its ID for a given user
// @Param		userId path int true "User ID"
// @Param		ccId path int true "ID of the Credit Card"
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/credit-cards/{ccId} [put]
func (fmh *FinanceManagerHandler) UpdateCreditCard(w http.ResponseWriter, r *http.Request) {
	method := "creditcard_handler.UpdateCreditCard"
	klogger.Enter(method)

	var payload models.CreditCard
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	ccId, err1 := strconv.Atoi(chi.URLParam(r, "ccId"))

	if err != nil {
		klogger.ExitError(method, constants.UserForbiddenToViewOtherUserDataError, err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		return
	}

	if err1 != nil {
		klogger.ExitError(method, "unexpected error occured when processing credit card id", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Read in credit card from payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, "failed to parse JSON object", err)
		return
	}

	// Validate that the credit card exists
	cc, err := fmh.DB.GetCreditCardByID(ccId)
	if err != nil || cc.ID == 0 {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusNotFound)
		klogger.ExitError(method, constants.EntityNotFoundError, err)
		return
	}

	// Validate that the credit card belongs to the user. Return Not Found error to mask existence
	err = fmh.Validator.CreditCardBelongsToUser(cc, userId)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.EntityNotFoundError), http.StatusNotFound)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	//Validate the credit card object
	payload.ValidateCanSaveCreditCard()

	// Update the credit card
	err = fmh.DB.UpdateCreditCard(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.FailedToUpdateEntityError, err)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, constants.SuccessMessage)
}
