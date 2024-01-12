package fmhandler

import (
	"finance-manager-backend/cmd/finance-mngr/internal/authentication"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/jsonutils"
	"finance-manager-backend/cmd/finance-mngr/internal/repository"
	"finance-manager-backend/cmd/finance-mngr/internal/stockservice.go"
	"finance-manager-backend/cmd/finance-mngr/internal/validation"
	"os"
)

type FinanceManagerHandler struct {
	JSONUtil             jsonutils.JSONUtils
	DB                   repository.DatabaseRepo
	Auth                 authentication.Auth
	Validator            validation.AppValidator
	Version              string
	PolygonApiKey        string
	StocksEnabled        bool
	StocksApiKeyFileName string
	StocksService        stockservice.StockService
}

func (fmh FinanceManagerHandler) GetVersion() string {
	return fmh.Version
}

// Attempts to load an API key from a file
func (fmh *FinanceManagerHandler) LoadApiKeyFromFile() error {
	method := "polygon_api.LoadApiKeyFromFile"
	fmlogger.Enter(method)

	pwd, _ := os.Getwd()
	bs, err := os.ReadFile(pwd + fmh.StocksApiKeyFileName)

	if err != nil {
		fmlogger.ExitError(method, "key file not found", err)
		return err
	}

	fmlogger.Info(method, "key loaded successfully")

	fmh.PolygonApiKey = string(bs)
	fmh.StocksEnabled = true

	fmlogger.Exit(method)
	return nil
}

// Reads an API key into the application object and persists it into a file
func (fmh *FinanceManagerHandler) UpdateAndPersistAPIKey(k string) error {
	method := "polygon_api.UpdateAndPersistAPIKey"
	fmlogger.Enter(method)

	fmlogger.Info(method, "Loading key into application")
	fmh.PolygonApiKey = k
	fmh.StocksEnabled = true

	fmlogger.Info(method, "attempting to persist API key file")
	pwd, _ := os.Getwd()

	err := os.WriteFile(pwd+fmh.StocksApiKeyFileName, []byte(k), 0666)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when writing key file", err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}
