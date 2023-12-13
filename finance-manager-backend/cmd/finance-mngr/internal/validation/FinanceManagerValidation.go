package validation

import (
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"finance-manager-backend/cmd/finance-mngr/internal/repository"
)

type AppValidator interface {

	//Users
	CheckIfUserHasRole(id int, desiredRole string) (bool, error)
	IsValidToEnterNewUser(user models.User) error
	IsValidToViewOtherUserData(loggedInUserId int) (bool, error)

	//Loans
	LoanBelongsToUser(loan models.Loan, userId int) error

	//Incomes
	IncomeBelongsToUser(income models.Income, userId int) error
}

type FinanceManagerValidator struct {
	DB repository.DatabaseRepo
}
