package service

import (
	"finance-manager-backend/internal/finance-mngr/models"
	"time"
)

type ExternalService interface {

	//Loads in API key for external stock calls
	UpdateAndPersistAPIKey(k string) error

	//Fetches a response indicating if stocks are enabled or not
	GetIsStocksEnabled() bool

	//Loads the stock API from a file
	LoadApiKeyFromFile() error

	//Fetches a Stock for a given ticker
	FetchStockWithTicker(ticker string) (models.Stock, error)

	//Fetches the past 1 year of data for a given ticker
	FetchStockWithTickerForPastYear(ticker string) ([]models.Stock, error)

	//Fetches stocks for a given ticker and date range
	FetchStockWithTickerForDateRange(t string, d1 time.Time, d2 time.Time) ([]models.Stock, error)
}
