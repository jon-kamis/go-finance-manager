package fmhandler

import (
	"finance-manager-backend/cmd/finance-mngr/internal/authentication"
	"finance-manager-backend/cmd/finance-mngr/internal/jsonutils"
	"finance-manager-backend/cmd/finance-mngr/internal/repository"
	"finance-manager-backend/cmd/finance-mngr/internal/stockservice.go"
	"finance-manager-backend/cmd/finance-mngr/internal/validation"
)

// @title Go Finance Manager API
// @version 1.0.0
// @description This API serves personal finance endpoints. Accuracy is not garunteed
// @BasePath /
type FinanceManagerHandler struct {
	JSONUtil      jsonutils.JSONUtils
	DB            repository.DatabaseRepo
	Auth          authentication.Auth
	Validator     validation.AppValidator
	Version       string
	StocksService stockservice.StockService
}

func (fmh FinanceManagerHandler) GetVersion() string {
	return fmh.Version
}
