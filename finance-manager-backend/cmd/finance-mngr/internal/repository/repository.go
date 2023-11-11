package repository

import (
	"database/sql"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByUsernameOrEmail(username string, email string) (*models.User, error)
	InsertUser(models.User) error
	UpdateUserDetails(models.User) error
	GetUserRoles(id int) ([]*models.UserRole, error)
}
