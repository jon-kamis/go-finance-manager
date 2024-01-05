package application

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"net/http"
)

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
