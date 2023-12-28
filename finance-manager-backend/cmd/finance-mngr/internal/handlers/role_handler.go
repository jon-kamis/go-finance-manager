package handlers

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"net/http"
)

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
