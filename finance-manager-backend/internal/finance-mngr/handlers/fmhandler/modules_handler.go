package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetIsModuleEnabled godoc
// @title		Module Enabled
// @version 	1.0.0
// @Tags 		Modules
// @Summary 	Module Enabled
// @Description Returns a boolean stating whether the requested module is enabled or not
// @Param		moduleName path string true "The name of the module to check. Options are {stocks}"
// @Produce 	json
// @Success 	200 {object} models.ModuleEnabledResponse
// @Failure		404 {object} jsonutils.JSONResponse
// @Router 		/modules/{moduleName} [get]
func (fmh *FinanceManagerHandler) GetIsModuleEnabled(w http.ResponseWriter, r *http.Request) {
	method := "modules_handler.GetIsModuleEnabled"
	fmlogger.Enter(method)

	name := chi.URLParam(r, "moduleName")

	re := models.ModuleEnabledResponse{}

	switch name {
	case constants.StockModuleName:
		re.Enabled = fmh.StocksService.GetIsStocksEnabled()
	default:
		//Requested module does not exist
		err := errors.New(constants.GenericNotFoundError)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusNotFound)
		fmlogger.ExitError(method, "module not found", err)
		return
	}

	fmh.JSONUtil.WriteJSON(w, http.StatusOK, re)
	fmlogger.Exit(method)
}

// PostModuleAPIKey godoc
// @title		Add Module API Key
// @version 	1.0.0
// @Tags 		Modules
// @Summary 	Add Module API key
// @Description Adds or overwrites the API key for the given module if allowed for this module
// @Param		moduleName path string true "The name of the module to add a key for. Options are {stocks}"
// @Param		keyRequest body models.EnableModuleRequest true "The request containing the Key to add"
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure		404 {object} jsonutils.JSONResponse
// @Router 		/modules/{moduleName}/key [post]
func (fmh *FinanceManagerHandler) PostModuleAPIKey(w http.ResponseWriter, r *http.Request) {
	method := "modules_handler.PostStocksAPIKey"
	fmlogger.Enter(method)

	uId, err := fmh.Auth.GetLoggedInUserId(w, r)

	//uId must be loaded successfully to proceed
	if err != nil {
		fmlogger.ExitError(method, constants.FailedToReadUserIdFromAuthHeaderError, err)
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.FailedToReadUserIdFromAuthHeaderError), http.StatusInternalServerError)
		return
	}

	//Determine if user has admin role
	hasRole, err := fmh.Validator.CheckIfUserHasRole(uId, "admin")

	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	//User must be admin to proceed
	if !hasRole {
		err = errors.New(constants.GenericForbiddenError)
		fmlogger.ExitError(method, err.Error(), err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusForbidden)
		return
	}

	var payload models.EnableModuleRequest

	// Read payload
	err = fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmlogger.ExitError(method, constants.GenericBadRequestError, err)
		fmh.JSONUtil.ErrorJSON(w, errors.New(constants.GenericBadRequestError), http.StatusBadRequest)
		return
	}

	// Determine the module we are trying to add this key for
	name := chi.URLParam(r, "moduleName")

	switch name {
	case constants.StockModuleName:
		err = fmh.StocksService.UpdateAndPersistAPIKey(payload.Key)
		if err != nil {
			fmlogger.ExitError(method, constants.GenericServerError, err)
			fmh.JSONUtil.ErrorJSON(w, errors.New(constants.GenericServerError), http.StatusInternalServerError)
			return
		}
	default:
		//Requested module does not exist
		err := errors.New(constants.GenericNotFoundError)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusNotFound)
		fmlogger.ExitError(method, "module not found", err)
		return
	}

	fmh.JSONUtil.WriteJSON(w, http.StatusOK, constants.SuccessMessage)
	fmlogger.Exit(method)
}
