package models

type StockPortfolioHistoryResponse struct {
	Items []PortfolioBalanceHistory `json:"items"`
	Count int                       `json:"count"`
}
