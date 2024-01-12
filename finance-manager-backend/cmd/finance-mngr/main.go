package main

import (
	"finance-manager-backend/cmd/finance-mngr/internal/application"
	"finance-manager-backend/cmd/finance-mngr/internal/authentication"
	"finance-manager-backend/cmd/finance-mngr/internal/config"
	"finance-manager-backend/cmd/finance-mngr/internal/constants"
	"finance-manager-backend/cmd/finance-mngr/internal/handlers/fmhandler"
	"finance-manager-backend/cmd/finance-mngr/internal/jsonutils"
	"finance-manager-backend/cmd/finance-mngr/internal/repository/dbrepo"
	"finance-manager-backend/cmd/finance-mngr/internal/stockservice.go/fmstockservice"
	"finance-manager-backend/cmd/finance-mngr/internal/validation"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const port = 8080

func main() {
	//set application config
	var app application.Application
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
	conn, err := app.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.Auth = authentication.Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 60,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "Host-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	app.JSONUtil = &jsonutils.JSONUtil{}

	stockService := fmstockservice.FmStockService{
		StocksEnabled:        false,
		StocksApiKeyFileName: constants.APIKeyFileName,
	}

	stockService.LoadApiKeyFromFile()

	app.StocksService = &stockService
	app.Handler = &fmhandler.FinanceManagerHandler{
		JSONUtil:  &jsonutils.JSONUtil{},
		DB:        app.DB,
		Auth:      app.Auth,
		Validator: &validation.FinanceManagerValidator{DB: app.DB},

		StocksService: &stockService,
	}

	defer app.DB.Connection().Close()

	app.Domain = "fm.com"

	log.Println("Starting application on port", port)

	//Kick off Scheduled Jobs
	go scheduledMinuteJobs(time.NewTicker(time.Minute*1), app)

	//start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.Routes())
	if err != nil {
		log.Fatal(err)
	}
}
