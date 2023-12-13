package repository

import (
	"database/sql"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB

	// User functions
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByUsernameOrEmail(username string, email string) (*models.User, error)
	GetAllUsers(search string) ([]*models.User, error)
	InsertUser(models.User) (int, error)
	UpdateUserDetails(models.User) error

	// Role functions
	GetRoleByCode(string) (*models.Role, error)

	// User Role functions
	GetUserRoles(id int) ([]*models.UserRole, error)
	InsertUserRole(models.UserRole) (int, error)

	// Loan Functions
	DeleteLoanByID(id int) error
	GetAllUserLoans(userId int, search string) ([]*models.Loan, error)
	GetLoanByID(id int) (models.Loan, error)
	InsertLoan(models.Loan) (int, error)
	UpdateLoan(loan models.Loan) error

	//Income Functions
	DeleteIncomeByID(id int) error
	GetAllUserIncomes(id int, search string) ([]*models.Income, error)
	GetIncomeByID(id int) (models.Income, error)
	InsertIncome(models.Income) (int, error)
	UpdateIncome(income models.Income) error
}
