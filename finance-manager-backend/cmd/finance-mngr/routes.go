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
	mux.Get("/roles", app.Handler.GetAllRoles)

	mux.Route("/users", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Get("/all", app.Handler.GetAllUsers)
		mux.Delete("/{userId}", app.Handler.DeleteUserById)
		mux.Get("/{userId}/roles", app.Handler.GetUserRoles)

		mux.Get("/{userId}", app.Handler.GetUserByID)
		mux.Get("/{userId}/summary", app.Handler.GetUserSummary)

		mux.Get("/{userId}/loans", app.Handler.GetAllUserLoans)
		mux.Post("/{userId}/loans", app.Handler.SaveLoan)
		mux.Get("/{userId}/loan-summary", app.Handler.GetLoanSummary)
		mux.Get("/{userId}/loans/{loanId}", app.Handler.GetLoanById)
		mux.Put("/{userId}/loans/{loanId}", app.Handler.UpdateLoan)
		mux.Delete("/{userId}/loans/{loanId}", app.Handler.DeleteLoanById)
		mux.Post("/{userId}/loans/{loanId}/calculate", app.Handler.CalculateLoan)
		mux.Post("/{userId}/loans/{loanId}/compare-payments", app.Handler.CompareLoanPayments)

		mux.Get("/{userId}/incomes", app.Handler.GetAllUserIncomes)
		mux.Get("/{userId}/incomes/{incomeId}", app.Handler.GetIncomeById)
		mux.Put("/{userId}/incomes/{incomeId}", app.Handler.UpdateIncome)
		mux.Delete("/{userId}/incomes/{incomeId}", app.Handler.DeleteIncomeById)
		mux.Post("/{userId}/incomes", app.Handler.SaveIncome)

		mux.Get("/{userId}/bills", app.Handler.GetAllUserBills)
		mux.Post("/{userId}/bills", app.Handler.SaveBill)
		mux.Get("/{userId}/bills/{billId}", app.Handler.GetBillById)
		mux.Put("/{userId}/bills/{billId}", app.Handler.UpdateBill)
		mux.Delete("/{userId}/bills/{billId}", app.Handler.DeleteBillById)
	})

	return mux
}
