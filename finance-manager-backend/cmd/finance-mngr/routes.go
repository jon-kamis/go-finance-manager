package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	// Create a router mux
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Handler.Home)

	mux.Post("/authenticate", app.Handler.Authenticate)
	mux.Get("/refresh", app.Handler.RefreshToken)
	mux.Get("/logout", app.Handler.Logout)
	mux.Post("/register", app.Handler.Register)
	mux.Route("/users", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Get("/all", app.Handler.GetAllUsers)
		mux.Get("/{userId}", app.Handler.GetUserByID)

		mux.Get("/{userId}/loans", app.Handler.GetAllUserLoans)
		mux.Post("/{userId}/loans", app.Handler.SaveLoan)
		mux.Get("/{userId}/loan-summary", app.Handler.GetLoanSummary)
		mux.Get("/{userId}/loans/{loanId}", app.Handler.GetLoanById)
		mux.Put("/{userId}/loans/{loanId}", app.Handler.UpdateLoan)
		mux.Delete("/{userId}/loans/{loanId}", app.Handler.DeleteLoanById)
		mux.Post("/{userId}/loans/{loanId}/calculate", app.Handler.CalculateLoan)
		mux.Post("/{userId}/loans/{loanId}/compare-payments", app.Handler.CompareLoanPayments)
	})

	return mux
}
