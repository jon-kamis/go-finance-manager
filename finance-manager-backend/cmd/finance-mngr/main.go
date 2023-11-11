package main

import (
	"finance-manager-backend/cmd/finance-mngr/internal/repository"
	"finance-manager-backend/cmd/finance-mngr/internal/repository/dbrepo"
	"fmt"
	"log"
	"net/http"
	"time"
)

const port = 8080

type application struct {
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	FrontendUrl  string
}

func main() {
	//set application config
	var app application
	config := FinanceManagerConfig{
		DSN: env_value{
			envName:    "DSN",
			defaultVal: "host=localhost port=5432 user=postgres password=postgres dbname=financemanager sslmode=disable timezone=UTC connect_timeout=5",
		},
		JWTSecret: env_value{
			envName:    "JWTSecret",
			defaultVal: "verysecret",
		},
		JWTIssuer: env_value{
			envName:    "JWTIssuer",
			defaultVal: "fm.com",
		},
		JWTAudience: env_value{
			envName:    "JWTAudience",
			defaultVal: "fm.com",
		},
		CookieDomain: env_value{
			envName:    "CookieDomain",
			defaultVal: "localhost",
		},
		Domain: env_value{
			envName:    "Domain",
			defaultVal: "fm.com",
		},
		FrontendUrl: env_value{
			envName:    "FrontendUrl",
			defaultVal: "http://localhost:3000",
		},
	}

	//Attempt to read values from the environment
	app.DSN = getEnvFromEnvValue(config.DSN)
	app.JWTSecret = getEnvFromEnvValue(config.JWTSecret)
	app.JWTIssuer = getEnvFromEnvValue(config.JWTIssuer)
	app.JWTAudience = getEnvFromEnvValue(config.JWTAudience)
	app.CookieDomain = getEnvFromEnvValue(config.CookieDomain)
	app.Domain = getEnvFromEnvValue(config.Domain)
	app.FrontendUrl = getEnvFromEnvValue(config.FrontendUrl)

	//connect to db
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.auth.Secret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "Host-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	app.Domain = "example.com"

	log.Println("Starting application on port", port)

	//start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
