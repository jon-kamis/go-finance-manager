package fmhandler

import (
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models/restmodels"
	"net/http"

	"github.com/jon-kamis/klogger"
)

// CalcSavingsRequest godoc
// @title		Calculate Savings
// @version 	1.0.0
// @Tags 		Savings
// @Summary 	Calculate Savings Request
// @Description Performs calculation on request and returns a result without saving
// @Param		request body restmodels.SavingsCalculationRequest true "the request to calculate"
// @Produce 	json
// @Success 	200 {object} models.Income
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/calc-savings [post]
func (fmh *FinanceManagerHandler) CalcSavingsRequest(w http.ResponseWriter, r *http.Request) {
	method := "savings_handler.CalcSavingsRequest"
	klogger.Enter(method)

	var payload restmodels.SavingsCalculationRequest

	// Read in loan from payload
	err := fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, constants.FailedToParseJsonBodyError, err)
		return
	}

	resp, err := payload.Calculate()

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, err.Error())
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, resp)
}
