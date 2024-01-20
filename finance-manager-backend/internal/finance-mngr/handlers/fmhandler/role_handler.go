package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"net/http"
)

// GetAllRoles godoc
// @title		Get All Roles
// @version 	1.0.0
// @Tags 		Roles
// @Summary 	Get All Roles
// @Description Returns an array of Role objects
// @Param		search query string false "Search for roles by name"
// @Produce 	json
// @Success 	200 {array} models.Role
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/roles [get]
func (fmh *FinanceManagerHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	method := "role_handler.GetAllRoles"
	fmlogger.Enter(method)

	search := r.URL.Query().Get("search")
	roles, err := fmh.DB.GetAllRoles(search)

	if err != nil {
		fmlogger.ExitError(method, "error occured when fetching roles", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching roles list"), http.StatusInternalServerError)
		return
	}

	fmlogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, roles)
}
