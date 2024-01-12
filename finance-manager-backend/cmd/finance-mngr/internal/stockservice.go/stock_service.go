package stockservice

import "finance-manager-backend/cmd/finance-mngr/internal/models"

type StockService interface {
	FetchStockWithTicker(ticker string, token string) (models.Stock, error)
}
