package validation

import (
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/internal/finance-mngr/repository"
)

type AppValidator interface {

	//Users
	CheckIfUserHasRole(id int, desiredRole string) (bool, error)
	IsValidToEnterNewUser(user models.User) error
	IsValidToViewOtherUserData(loggedInUserId int) (bool, error)
	IsValidToDeleteOtherUserData(loggedInUserId int) (bool, error)

	//UserRoles
	UserRoleExistsAndBelongsToUser(roleId, userId int) error
	UserRoleBelongsToUser(userRole models.UserRole, userId int) error

	//Loans
	LoanBelongsToUser(loan models.Loan, userId int) error

	//Incomes
	IncomeBelongsToUser(income models.Income, userId int) error

	//Bills
	BillBelongsToUser(bill models.Bill, userId int) error

	//Credit Cards
	CreditCardBelongsToUser(cc models.CreditCard, userId int) error
}

type FinanceManagerValidator struct {
	DB repository.DatabaseRepo
}
