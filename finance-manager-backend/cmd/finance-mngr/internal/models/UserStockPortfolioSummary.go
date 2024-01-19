package models

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"math"
	"time"
)

type UserStockPortfolioSummary struct {
	CurrentValue float64             `json:"currentValue"`
	CurrentHigh  float64             `json:"currentHigh"`
	CurrentLow   float64             `json:"currentLow"`
	CurrentOpen  float64             `json:"currentOpen"`
	CurrentClose float64             `json:"currentClose"`
	AsOfDate     time.Time           `json:"asOf"`
	Positions    []PortfolioPosition `json:"positions"`
}

type PortfolioPosition struct {
	Ticker   string    `json:"ticker"`
	Quantity float64   `json:"quantity"`
	Value    float64   `json:"value"`
	Open     float64   `json:"open"`
	Close    float64   `json:"close"`
	High     float64   `json:"high"`
	Low      float64   `json:"low"`
	AsOfDate time.Time `json:"asOf"`
}

// Calculates the daily totals for a portfolio. Expects Positions to be loaded
func (u *UserStockPortfolioSummary) CalcDailyTotals() {
	method := "UserStockPortfolioSummary.CalcDailyTotals"
	fmlogger.Enter(method)

	var h float64   //high
	var l float64   //low
	var o float64   //open
	var c float64   //close
	var t float64   //total
	d := time.Now() //as of date

	for _, p := range u.Positions {
		h += (p.High * p.Quantity)
		l += (p.Low * p.Quantity)
		o += (p.Open * p.Quantity)
		c += (p.Close * p.Quantity)
		t += p.Value

		if p.AsOfDate.Before(d) {
			d = p.AsOfDate
		}
	}

	u.CurrentHigh = math.Round(h*100) / 100
	u.CurrentLow = math.Round(l*100) / 100
	u.CurrentOpen = math.Round(o*100) / 100
	u.CurrentClose = math.Round(c*100) / 100
	u.CurrentValue = math.Round(t*100) / 100
	u.AsOfDate = d

	fmlogger.Exit(method)
}

// Loads Positions into the summary and triggers the calculation of daily total values
func (u *UserStockPortfolioSummary) LoadPositions(p []PortfolioPosition) {
	method := "UserStockPortfolioSummary.LoadPositions"
	fmlogger.Enter(method)

	u.Positions = p
	u.CalcDailyTotals()

	fmlogger.Exit(method)
}
