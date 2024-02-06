package fmhandler

import (
	"finance-manager-backend/internal/finance-mngr/models/restmodels"
	"html/template"
	"net/http"

	"github.com/jon-kamis/klogger"
)

// Home godoc
// @title Home
// @version 1.0.0
// @Tags Home
// @Summary Home
// @Description Returns application information and health
// @Produce json
// @Success 200 {object} restmodels.HomeResponse
// @Router / [get]
func (fmh *FinanceManagerHandler) Home(w http.ResponseWriter, r *http.Request) {

	method := "home_handler.Home"
	klogger.Enter(method)

	tmpl := template.Must(template.ParseFiles("./web/template/home.html"))

	var apiInfo = restmodels.HomeResponse{
		Status:  "active",
		Message: "Finance Manager Backend up and running!",
		Version: fmh.GetVersion(),
	}

	tmpl.Execute(w, struct {
		ApiInfo restmodels.HomeResponse
		Port    int
	}{
		ApiInfo: apiInfo,
		Port:    fmh.ApiPort,
	})
	klogger.Exit(method)
}
