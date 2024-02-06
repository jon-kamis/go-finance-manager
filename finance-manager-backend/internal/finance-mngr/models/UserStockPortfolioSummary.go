package models

import (
	"math"
	"time"

	"github.com/jon-kamis/klogger"
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

// Type PortfolioBalanceHistory Holds a record for the overall balance of the user's stocks for a given date
type PortfolioBalanceHistory struct {
	Date  time.Time `json:"date"`
	High  float64   `json:"high"`
	Low   float64   `json:"low"`
	Open  float64   `json:"open"`
	Close float64   `json:"close"`
}

// Type PortfolioPosition holds values for a user's Stock
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

// Type PositionHistory holds historic values for a Stock
type PositionHistory struct {
	Ticker          string    `json:"ticker"`
	High            float64   `json:"high"`
	Low             float64   `json:"low"`
	Open            float64   `json:"open"`
	Close           float64   `json:"close"`
	Delta           float64   `json:"delta"`
	DeltaPercentage float64   `json:"deltaPercentage"`
	Count           int       `json:"count"`
	StartDt         time.Time `json:"startDt"`
	EndDt           time.Time `json:"endDt"`
	Values          []Stock   `json:"values"`
}

// Calculates the daily totals for a portfolio. Expects Positions to be loaded
func (u *UserStockPortfolioSummary) CalcDailyTotals() {
	method := "UserStockPortfolioSummary.CalcDailyTotals"
	klogger.Enter(method)

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

	klogger.Exit(method)
}

// Loads Positions into the summary and triggers the calculation of daily total values
func (u *UserStockPortfolioSummary) LoadPositions(p []PortfolioPosition) {
	method := "UserStockPortfolioSummary.LoadPositions"
	klogger.Enter(method)

	u.Positions = p
	u.CalcDailyTotals()

	klogger.Exit(method)
}
