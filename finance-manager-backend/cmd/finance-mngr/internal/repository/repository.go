package repository

import (
	"database/sql"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB

	/*** User functions ***/

	//Deletes a user object by its id
	DeleteUserByID(id int) error

	//Fetches a user object by its id
	GetUserByID(id int) (*models.User, error)

	//Fetches a user object by its username
	GetUserByUsername(username string) (*models.User, error)

	//Fetches a user object by either its username or its email
	GetUserByUsernameOrEmail(username string, email string) (*models.User, error)

	//Fetches all user objects
	GetAllUsers(search string) ([]*models.User, error)

	//Inserts a user object
	InsertUser(models.User) (int, error)

	//Updates an existing user object by its id
	UpdateUserDetails(models.User) error

	/*** Role functions ***/

	//Fetches a role object by its code
	GetRoleByCode(string) (*models.Role, error)

	//Fetches a role object by its id
	GetRoleById(string) (*models.Role, error)

	//Fetches all role objects
	GetAllRoles(string) ([]*models.Role, error)

	/*** User Role functions ***/

	//Deletes User roles by their userId
	DeleteUserRolesByUserID(id int) error

	//Deletes a User role by its id
	DeleteUserRoleByID(id int) error

	//Fetches all user roles for a given userId
	GetAllUserRoles(id int) ([]*models.UserRole, error)

	//Fetches a User Role by its id
	GetUserRoleByID(userRoleId int) (models.UserRole, error)

	//Fetches a User Role by its roleId and userId
	GetUserRoleByRoleIDAndUserID(roleId int, uId int) (models.UserRole, error)

	//Inserts a new User Role
	InsertUserRole(models.UserRole) (int, error)

	/*** Loan Functions ***/

	//Deletes Loans by their userId
	DeleteLoansByUserID(id int) error

	//Delets a Loan by its id
	DeleteLoanByID(id int) error

	//Fetches all loans for a given userId and accepts a search parameter
	GetAllUserLoans(userId int, search string) ([]*models.Loan, error)

	//Fetches a Loan by its id
	GetLoanByID(id int) (models.Loan, error)

	//Inserts a new Loan
	InsertLoan(models.Loan) (int, error)

	//Updates an existing loan
	UpdateLoan(loan models.Loan) error

	//Income Functions
	DeleteIncomesByUserID(id int) error
	DeleteIncomeByID(id int) error
	GetAllUserIncomes(id int, search string) ([]*models.Income, error)
	GetIncomeByID(id int) (models.Income, error)
	InsertIncome(models.Income) (int, error)
	UpdateIncome(income models.Income) error

	//Bill Functions
	DeleteBillsByUserID(id int) error
	DeleteBillByID(id int) error
	GetAllUserBills(id int, search string) ([]*models.Bill, error)
	GetBillByID(id int) (models.Bill, error)
	InsertBill(models.Bill) (int, error)
	UpdateBill(income models.Bill) error

	//Credit Cards
	GetAllUserCreditCards(id int, search string) ([]*models.CreditCard, error)
	GetCreditCardByID(id int) (models.CreditCard, error)
	DeleteCreditCardsByID(id int) error
	DeleteCreditCardsByUserID(id int) error
	InsertCreditCard(cc models.CreditCard) (int, error)
	UpdateCreditCard(cc models.CreditCard) error

	/*** Stocks ***/

	//Inserts a new stock object
	InsertStock(s models.Stock) (int, error)

	//Gets a Stock by its ticker
	GetStockByTicker(t string) (models.Stock, error)

	//Fetches the stock that has both the oldest date and last_update_dt
	GetOldestStock() (models.Stock, error)

	UpdateStock(s models.Stock) error

	/*** Stock Data ***/

	//Inserts Stock Data
	InsertStockData(sl []models.Stock) error

	//Fetches latest stock data for a given ticker
	GetLatestStockDataByTicker(t string) (models.Stock, error)

	/*** User Stocks ***/

	//Inserts a new user stock object
	InsertUserStock(s models.UserStock) (int, error)

	//Fetches all UserStocks for a given user and accepts a search string
	GetAllUserStocks(userId int, search string) ([]*models.UserStock, error)
}
