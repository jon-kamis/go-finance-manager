package fmhandler

import (
	"finance-manager-backend/internal/finance-mngr/authentication"
	"finance-manager-backend/internal/finance-mngr/jsonutils"
	"finance-manager-backend/internal/finance-mngr/repository"
	"finance-manager-backend/internal/finance-mngr/service"
	"finance-manager-backend/internal/finance-mngr/stockservice"
	"finance-manager-backend/internal/finance-mngr/validation"
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
	Service       service.Service
	StocksService stockservice.StockService
	ApiPort       int
}

func (fmh FinanceManagerHandler) GetVersion() string {
	return fmh.Version
}
