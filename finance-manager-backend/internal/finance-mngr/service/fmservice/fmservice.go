package fmservice

import (
	"finance-manager-backend/internal/finance-mngr/repository"
)

type FMService struct {
	DB repository.DatabaseRepo
}
