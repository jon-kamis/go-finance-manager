package handlers

import (
	"net/http"
)

func (fmh *FinanceManagerHandler) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Finance Manager Backend up and running!",
		Version: "1.0.0",
	}

	_ = fmh.JSONUtil.WriteJSON(w, http.StatusOK, payload)
}
