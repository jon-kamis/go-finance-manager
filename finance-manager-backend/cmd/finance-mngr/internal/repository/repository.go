package repository

import (
	"database/sql"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB

	// User functions
	DeleteUserByID(id int) error
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByUsernameOrEmail(username string, email string) (*models.User, error)
	GetAllUsers(search string) ([]*models.User, error)
	InsertUser(models.User) (int, error)
	UpdateUserDetails(models.User) error

	// Role functions
	GetRoleByCode(string) (*models.Role, error)
	GetRoleById(string) (*models.Role, error)
	GetAllRoles(string) ([]*models.Role, error)

	// User Role functions
	DeleteUserRolesByUserID(id int) error
	DeleteUserRoleByID(id int) error
	GetAllUserRoles(id int) ([]*models.UserRole, error)
	GetUserRoleByID(userRoleId int) (models.UserRole, error)
	GetUserRoleByRoleIDAndUserID(roleId int, uId int) (models.UserRole, error)
	InsertUserRole(models.UserRole) (int, error)

	// Loan Functions
	DeleteLoansByUserID(id int) error
	DeleteLoanByID(id int) error
	GetAllUserLoans(userId int, search string) ([]*models.Loan, error)
	GetLoanByID(id int) (models.Loan, error)
	InsertLoan(models.Loan) (int, error)
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
}
