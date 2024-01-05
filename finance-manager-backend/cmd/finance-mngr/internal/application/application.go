package application

import (
	"finance-manager-backend/cmd/finance-mngr/internal/authentication"
	"finance-manager-backend/cmd/finance-mngr/internal/handlers"
	"finance-manager-backend/cmd/finance-mngr/internal/jsonutils"
	"finance-manager-backend/cmd/finance-mngr/internal/repository"
)

type Application struct {
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	Auth         authentication.Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	FrontendUrl  string
	Handler      handlers.Handler
	JSONUtil     jsonutils.JSONUtils
}
