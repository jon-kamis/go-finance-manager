package fmhandler

import (
	models "finance-manager-backend/internal/finance-mngr/models/rest"
	"net/http"
)

// Home godoc
// @title Home
// @version 1.0.0
// @Tags Home
// @Summary Home
// @Description Returns application information and health
// @Produce json
// @Success 200 {object} models.HomeResponse
// @Router / [get]
func (fmh *FinanceManagerHandler) Home(w http.ResponseWriter, r *http.Request) {
	var payload = models.HomeResponse{
		Status:  "active",
		Message: "Finance Manager Backend up and running!",
		Version: fmh.GetVersion(),
	}

	_ = fmh.JSONUtil.WriteJSON(w, http.StatusOK, payload)
}
