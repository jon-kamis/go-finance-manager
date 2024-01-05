package fmhandler

import (
	"finance-manager-backend/cmd/finance-mngr/internal/authentication"
	"finance-manager-backend/cmd/finance-mngr/internal/jsonutils"
	"finance-manager-backend/cmd/finance-mngr/internal/repository"
	"finance-manager-backend/cmd/finance-mngr/internal/validation"
)

type FinanceManagerHandler struct {
	JSONUtil  jsonutils.JSONUtils
	DB        repository.DatabaseRepo
	Auth      authentication.Auth
	Validator validation.AppValidator
	Version   string
}

func (fmh FinanceManagerHandler) GetVersion() string {
	return fmh.Version
}
