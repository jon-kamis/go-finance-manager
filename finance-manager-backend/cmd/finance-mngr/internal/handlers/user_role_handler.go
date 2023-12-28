package handlers

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (fmh *FinanceManagerHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	method := "user_role_handler.GetUserRoles"
	fmlogger.Enter(method)

	//Read ID from url
	search := r.URL.Query().Get("search")
	id, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), false, w, r)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching user roles", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	loans, err := fmh.DB.GetAllUserLoans(id, search)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching loans", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching user roles"), http.StatusInternalServerError)
		return
	}

	fmlogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, loans)
}
