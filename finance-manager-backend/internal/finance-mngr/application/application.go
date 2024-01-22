// Package application contains core files required to run a GO web API
package application

import (
	"finance-manager-backend/internal/finance-mngr/authentication"
	"finance-manager-backend/internal/finance-mngr/handlers"
	"finance-manager-backend/internal/finance-mngr/jsonutils"
	"finance-manager-backend/internal/finance-mngr/repository"
	"finance-manager-backend/internal/finance-mngr/stockservice"
)

// Type Application stores environment variables and objects required to run Finance Manager
type Application struct {
	DSN           string
	Domain        string
	DB            repository.DatabaseRepo
	Auth          authentication.Auth
	JWTSecret     string
	JWTIssuer     string
	JWTAudience   string
	CookieDomain  string
	FrontendUrl   string
	Handler       handlers.Handler
	JSONUtil      jsonutils.JSONUtils
	StocksService stockservice.StockService
}
