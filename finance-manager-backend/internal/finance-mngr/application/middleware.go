package application

import (
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"net/http"
)

//Function EnableCORS enables CORS security on API requests
func (app *Application) EnableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", app.FrontendUrl)

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-TOKEN, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

//Function AuthRequired is used as middleware for a router request and will block any request using it that does not have a valid
//JWT token issues by this application
func (app *Application) AuthRequired(next http.Handler) http.Handler {
	method := "middleware.AuthRequired"
	fmlogger.Enter(method)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, err := app.Auth.GetTokenFromHeaderAndVerify(w, r)
		if err != nil {
			fmlogger.ExitError(method, "unauthorized", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			fmlogger.Exit(method)
			next.ServeHTTP(w, r)
		}
	})
}
