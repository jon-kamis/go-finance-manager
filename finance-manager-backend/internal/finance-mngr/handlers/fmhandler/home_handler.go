package fmhandler

import (
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	models "finance-manager-backend/internal/finance-mngr/models/rest"
	"html/template"
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

	method := "home_handler.Home"
	fmlogger.Enter(method)

	tmpl := template.Must(template.ParseFiles("./web/template/home.html"))

	var apiInfo = models.HomeResponse{
		Status:  "active",
		Message: "Finance Manager Backend up and running!",
		Version: fmh.GetVersion(),
	}

	tmpl.Execute(w, struct {
		ApiInfo models.HomeResponse
		Port    int
	}{
		ApiInfo: apiInfo,
		Port: fmh.ApiPort,
	})
	fmlogger.Exit(method)
}
