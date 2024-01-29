package models

type StockPortfolioHistoryResponse struct {
	High            float64                   `json:"high"`
	Low             float64                   `json:"low"`
	Open            float64                   `json:"open"`
	Close           float64                   `json:"close"`
	Delta           float64                   `json:"delta"`
	DeltaPercentage float64                   `json:"deltaPercentage"`
	Items           []PortfolioBalanceHistory `json:"items"`
	Count           int                       `json:"count"`
}
