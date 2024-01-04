package main

import (
	"finance-manager-backend/cmd/finance-mngr/internal/authentication"
	"finance-manager-backend/cmd/finance-mngr/internal/config"
	"finance-manager-backend/cmd/finance-mngr/internal/handlers"
	"finance-manager-backend/cmd/finance-mngr/internal/jsonutils"
	"finance-manager-backend/cmd/finance-mngr/internal/repository"
	"finance-manager-backend/cmd/finance-mngr/internal/repository/dbrepo"
	"finance-manager-backend/cmd/finance-mngr/internal/validation"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const port = 8080

type application struct {
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	auth         authentication.Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	FrontendUrl  string
	Handler      handlers.FinanceManagerHandler
	JSONUtil     jsonutils.JSONUtils
}

func main() {
	//set application config
	var app application
	var appConfig = config.GetDefaultConfig()

	//Attempt to read values from the environment
	os.Setenv("TZ", config.GetEnvFromEnvValue(appConfig.TimeZone))
	app.DSN = config.GetEnvFromEnvValue(appConfig.DSN)
	app.JWTSecret = config.GetEnvFromEnvValue(appConfig.JWTSecret)
	app.JWTIssuer = config.GetEnvFromEnvValue(appConfig.JWTIssuer)
	app.JWTAudience = config.GetEnvFromEnvValue(appConfig.JWTAudience)
	app.CookieDomain = config.GetEnvFromEnvValue(appConfig.CookieDomain)
	app.Domain = config.GetEnvFromEnvValue(appConfig.Domain)
	app.FrontendUrl = config.GetEnvFromEnvValue(appConfig.FrontendUrl)

	//connect to db
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.auth = authentication.Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.auth.Secret,
		TokenExpiry:   time.Minute * 60,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "Host-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	app.JSONUtil = &jsonutils.JSONUtil{}
	app.Handler = handlers.FinanceManagerHandler{
		JSONUtil:  &jsonutils.JSONUtil{},
		DB:        app.DB,
		Auth:      app.auth,
		Validator: &validation.FinanceManagerValidator{DB: app.DB},
	}

	defer app.DB.Connection().Close()

	app.Domain = "example.com"

	log.Println("Starting application on port", port)

	//start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
