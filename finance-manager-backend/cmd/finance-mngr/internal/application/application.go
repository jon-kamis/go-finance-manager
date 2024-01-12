// Package application contains core files required to run a GO web API
package application

import (
	"finance-manager-backend/cmd/finance-mngr/internal/authentication"
	"finance-manager-backend/cmd/finance-mngr/internal/handlers"
	"finance-manager-backend/cmd/finance-mngr/internal/jsonutils"
	"finance-manager-backend/cmd/finance-mngr/internal/repository"
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
}
