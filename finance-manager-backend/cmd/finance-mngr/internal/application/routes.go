package application

import (
	"net/http"

	_ "finance-manager-backend/cmd/finance-mngr/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *Application) Routes() http.Handler {
	// Create a router r
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(app.EnableCORS)

	//Swagger Docs
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	r.Get("/", app.Handler.Home)

	r.Post("/authenticate", app.Handler.Authenticate)
	r.Get("/refresh", app.Handler.RefreshToken)
	r.Get("/logout", app.Handler.Logout)
	r.Post("/register", app.Handler.Register)

	r.Route("/roles", func(r chi.Router) {
		r.Use(app.AuthRequired)
		r.Get("/", app.Handler.GetAllRoles)
	})

	r.Route("/modules", func(r chi.Router) {
		r.Use(app.AuthRequired)

		r.Route("/{moduleName}", func(r chi.Router) {
			r.Get("/", app.Handler.GetIsModuleEnabled)
			r.Post("/key", app.Handler.PostModuleAPIKey)
		})

	})

	r.Route("/users", func(r chi.Router) {
		r.Use(app.AuthRequired)

		r.Get("/", app.Handler.GetAllUsers)

		r.Route("/{userId}", func(r chi.Router) {

			r.Delete("/", app.Handler.DeleteUserById)
			r.Get("/", app.Handler.GetUserByID)
			r.Get("/summary", app.Handler.GetUserSummary)

			//User Role Routes
			r.Route("/roles", func(r chi.Router) {
				r.Get("/", app.Handler.GetUserRoles)
				r.Post("/{roleId}", app.Handler.AddUserRoles)
				r.Delete("/{roleId}", app.Handler.DeleteUserRoles)
			})

			//Loans Routes
			r.Get("/loans-summary", app.Handler.GetLoanSummary)
			r.Route("/loans", func(r chi.Router) {
				r.Get("/", app.Handler.GetAllUserLoans)
				r.Post("/", app.Handler.SaveLoan)

				r.Route("/{loanId}", func(r chi.Router) {
					r.Get("/", app.Handler.GetLoanById)
					r.Put("/", app.Handler.UpdateLoan)
					r.Delete("/", app.Handler.DeleteLoanById)
					r.Post("/calculate", app.Handler.CalculateLoan)
					r.Post("/compare-payments", app.Handler.CompareLoanPayments)
				})

			})

			//Incomes Routes
			r.Route("/incomes", func(r chi.Router) {
				r.Get("/", app.Handler.GetAllUserIncomes)
				r.Post("/", app.Handler.SaveIncome)

				r.Route("/{incomeId}", func(r chi.Router) {
					r.Get("/", app.Handler.GetIncomeById)
					r.Put("/", app.Handler.UpdateIncome)
					r.Delete("/", app.Handler.DeleteIncomeById)
				})

			})

			//Bills Routes
			r.Route("/bills", func(r chi.Router) {
				r.Get("/", app.Handler.GetAllUserBills)
				r.Post("/", app.Handler.SaveBill)

				r.Route("/{billId}", func(r chi.Router) {
					r.Get("/", app.Handler.GetBillById)
					r.Put("/", app.Handler.UpdateBill)
					r.Delete("/", app.Handler.DeleteBillById)
				})
			})

			//Credit Card Routes
			r.Route("/credit-cards", func(r chi.Router) {
				r.Get("/", app.Handler.GetAllUserCreditCards)
				r.Post("/", app.Handler.SaveCreditCard)

				r.Route("/{ccId}", func(r chi.Router) {
					r.Get("/", app.Handler.GetCreditCardById)
					r.Delete("/", app.Handler.DeleteCreditCardById)
					r.Put("/", app.Handler.UpdateCreditCard)
				})
			})

			//Stocks
			r.Route("/stocks", func(r chi.Router) {
				r.Post("/", app.Handler.SaveUserStock)
				r.Get("/", app.Handler.GetUserStockPortfolioSummary)

			})
		})

	})

	return r
}
