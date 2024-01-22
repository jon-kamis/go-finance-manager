package handlers

import "net/http"

type Handler interface {

	/*** Bills ***/

	//Deletes a specific bill object by its id for a given user
	DeleteBillById(w http.ResponseWriter, r *http.Request)

	//Fetches all bills for a given user and accepts a search parameter
	GetAllUserBills(w http.ResponseWriter, r *http.Request)

	//Fetches a specific Bill by its id for a given user
	GetBillById(w http.ResponseWriter, r *http.Request)

	//Inserts a new Bill into the database for a given user
	SaveBill(w http.ResponseWriter, r *http.Request)

	//Updates a specific bill by its id for a given user
	UpdateBill(w http.ResponseWriter, r *http.Request)

	/*** Credit Cards ***/

	//Deletes a specific CreditCard by its id for a given user
	DeleteCreditCardById(w http.ResponseWriter, r *http.Request)

	//Fetches all CreditCards for a given user and accepts a search parameter
	GetAllUserCreditCards(w http.ResponseWriter, r *http.Request)

	//Fetches a specific CreditCard for a given user
	GetCreditCardById(w http.ResponseWriter, r *http.Request)

	//Inserts a new CreditCard into the database for a given user
	SaveCreditCard(w http.ResponseWriter, r *http.Request)

	//Updates a specific CreditCard by its id for a given user
	UpdateCreditCard(w http.ResponseWriter, r *http.Request)

	/*** Home ***/

	//Returns API information as a heartbeat
	Home(w http.ResponseWriter, r *http.Request)

	/*** Incomes ***/

	//Deletes a specific Income by its id for a given user
	DeleteIncomeById(w http.ResponseWriter, r *http.Request)

	//Fetches all Incomes for a given user and accepts a search parameter
	GetAllUserIncomes(w http.ResponseWriter, r *http.Request)

	//Fetches a specific Income by its id for a given user
	GetIncomeById(w http.ResponseWriter, r *http.Request)

	//Inserts a new Income into the database for a given user
	SaveIncome(w http.ResponseWriter, r *http.Request)

	//Updates a specific Income by its id for a given user
	UpdateIncome(w http.ResponseWriter, r *http.Request)

	/*** Loans ***/

	//Performs a payment schedule calculation on a Loan object
	CalculateLoan(w http.ResponseWriter, r *http.Request)

	//Compares the payment schedules of two Loan objects and returns a summary
	CompareLoanPayments(w http.ResponseWriter, r *http.Request)

	//Deletes a specific Loan by its id for a given user
	DeleteLoanById(w http.ResponseWriter, r *http.Request)

	//Fetches all Loans for a given user and accepts a search parameter
	GetAllUserLoans(w http.ResponseWriter, r *http.Request)

	//Fetches a specific Loan by its id for a given user
	GetLoanById(w http.ResponseWriter, r *http.Request)

	//Generates a summary from all Loans for a specific user
	GetLoanSummary(w http.ResponseWriter, r *http.Request)

	//Inserts a new Loan into the database for a given user
	SaveLoan(w http.ResponseWriter, r *http.Request)

	//Updates a specific Loan by its id for a given user
	UpdateLoan(w http.ResponseWriter, r *http.Request)

	/*** Login ***/

	//Validates supplied credentials then generates and returns a JWT TokenPair
	Authenticate(w http.ResponseWriter, r *http.Request)

	//Returns an expired refresh token cookie in an API response
	Logout(w http.ResponseWriter, r *http.Request)

	//Accepts a refresh token and generates a new JWT TokenPair with refreshed expiration date
	RefreshToken(w http.ResponseWriter, r *http.Request)

	/*** Modules ***/

	//Fetches a response indicating if a module is enabled or not
	GetIsModuleEnabled(w http.ResponseWriter, r *http.Request)

	//Initializes or overwrites the API key for a module
	PostModuleAPIKey(w http.ResponseWriter, r *http.Request)

	/*** Registration ***/

	//Validates and Inserts a new User into the database
	Register(w http.ResponseWriter, r *http.Request)

	/*** Roles ***/

	//Fetches all Role objects
	GetAllRoles(w http.ResponseWriter, r *http.Request)

	/*** Stocks ***/
	GetStockHistory(w http.ResponseWriter, r *http.Request)

	GetUserStockPortfolioHistory(w http.ResponseWriter, r *http.Request)

	/*** User Stocks ***/

	//Saves New User Stocks object
	SaveUserStock(w http.ResponseWriter, r *http.Request)

	//Gets a summary of a user's stock portfolio
	GetUserStockPortfolioSummary(w http.ResponseWriter, r *http.Request)

	/*** Users ***/

	//Deletes a specific user by id
	DeleteUserById(w http.ResponseWriter, r *http.Request)

	//Fetches all user objects
	GetAllUsers(w http.ResponseWriter, r *http.Request)

	//Fetches a specific user object by id
	GetUserByID(w http.ResponseWriter, r *http.Request)

	//Fetches a summary for a given user by id
	GetUserSummary(w http.ResponseWriter, r *http.Request)

	/** User Roles **/

	//Inserts a new UserRole into the database, granting access to a user
	AddUserRoles(w http.ResponseWriter, r *http.Request)

	//Deletes a UserRole, revoking access from a User
	DeleteUserRoles(w http.ResponseWriter, r *http.Request)

	//Fetches a list of user Roles for a given user
	GetUserRoles(w http.ResponseWriter, r *http.Request)

	/** Version **/

	//Fetches the current API version and returns it
	GetVersion() string
}
