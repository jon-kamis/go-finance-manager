package main

import (
	"finance-manager-backend/internal/finance-mngr/application"
	"finance-manager-backend/internal/finance-mngr/authentication"
	"finance-manager-backend/internal/finance-mngr/config"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/handlers/fmhandler"
	"finance-manager-backend/internal/finance-mngr/jobs"
	"finance-manager-backend/internal/finance-mngr/jsonutils"
	"finance-manager-backend/internal/finance-mngr/repository/dbrepo"
	"finance-manager-backend/internal/finance-mngr/service/fmservice"
	"finance-manager-backend/internal/finance-mngr/stockservice/fmstockservice"
	"finance-manager-backend/internal/finance-mngr/validation"
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
		DB:                   app.DB,
	}

	stockService.LoadApiKeyFromFile()

	app.StocksService = &stockService
	app.Handler = &fmhandler.FinanceManagerHandler{
		JSONUtil:      &jsonutils.JSONUtil{},
		DB:            app.DB,
		Auth:          app.Auth,
		Validator:     &validation.FinanceManagerValidator{DB: app.DB},
		Version:       constants.AppVersion,
		StocksService: &stockService,
		Service:       &fmservice.FMService{DB: app.DB},
		ApiPort:       port,
	}

	defer app.DB.Connection().Close()

	app.Domain = "fm.com"

	log.Println("Starting application on port", port)

	//Kick off Scheduled Jobs
	go jobs.ScheduledMinuteJobs(time.NewTicker(time.Minute*1), app)

	//start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.Routes())
	if err != nil {
		log.Fatal(err)
	}
}
