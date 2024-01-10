package handlers

import "net/http"

type Handler interface {

	//Bills
	DeleteBillById(w http.ResponseWriter, r *http.Request)
	GetAllUserBills(w http.ResponseWriter, r *http.Request)
	GetBillById(w http.ResponseWriter, r *http.Request)
	SaveBill(w http.ResponseWriter, r *http.Request)
	UpdateBill(w http.ResponseWriter, r *http.Request)

	//Credit Cards
	DeleteCreditCardById(w http.ResponseWriter, r *http.Request)
	GetAllUserCreditCards(w http.ResponseWriter, r *http.Request)
	GetCreditCardById(w http.ResponseWriter, r *http.Request)
	SaveCreditCard(w http.ResponseWriter, r *http.Request)
	UpdateCreditCard(w http.ResponseWriter, r *http.Request)

	//Home
	Home(w http.ResponseWriter, r *http.Request)

	//Incomes
	DeleteIncomeById(w http.ResponseWriter, r *http.Request)
	GetAllUserIncomes(w http.ResponseWriter, r *http.Request)
	GetIncomeById(w http.ResponseWriter, r *http.Request)
	SaveIncome(w http.ResponseWriter, r *http.Request)
	UpdateIncome(w http.ResponseWriter, r *http.Request)

	//Loans
	CalculateLoan(w http.ResponseWriter, r *http.Request)
	CompareLoanPayments(w http.ResponseWriter, r *http.Request)
	DeleteLoanById(w http.ResponseWriter, r *http.Request)
	GetAllUserLoans(w http.ResponseWriter, r *http.Request)
	GetLoanById(w http.ResponseWriter, r *http.Request)
	GetLoanSummary(w http.ResponseWriter, r *http.Request)
	SaveLoan(w http.ResponseWriter, r *http.Request)
	UpdateLoan(w http.ResponseWriter, r *http.Request)

	//Login
	Authenticate(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)

	//Registration
	Register(w http.ResponseWriter, r *http.Request)

	//Roles
	GetAllRoles(w http.ResponseWriter, r *http.Request)

	//Users
	DeleteUserById(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
	GetUserSummary(w http.ResponseWriter, r *http.Request)

	//User Roles
	AddUserRoles(w http.ResponseWriter, r *http.Request)
	DeleteUserRoles(w http.ResponseWriter, r *http.Request)
	GetUserRoles(w http.ResponseWriter, r *http.Request)

	//Version
	GetVersion() string
}
