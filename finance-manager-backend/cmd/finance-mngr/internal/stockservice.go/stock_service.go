package stockservice

import "finance-manager-backend/cmd/finance-mngr/internal/models"

//Interface StockService is used to handle external calls
type StockService interface {

	//Fetches a Stock for a given ticker
	FetchStockWithTicker(ticker string) (models.Stock, error)

	//Fetches the past 1 year of data for a given ticker
	FetchStockWithTickerForPastYear(ticker string) ([]models.Stock, error)

	//Loads in API key for external stock calls
	UpdateAndPersistAPIKey(k string) error

	//Fetches a response indicating if stocks are enabled or not
	GetIsStocksEnabled() bool

	//Loads the stock API from a file
	LoadApiKeyFromFile() error
}
