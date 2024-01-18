package models

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalcDailyTotals(t *testing.T) {
	method := "UserStockPortfolioSummary_test.TestCalcDailyTotals"
	fmlogger.Enter(method)

	var sum UserStockPortfolioSummary
	var pl []PortfolioPosition
	d1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	d3 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	pl = append(pl, PortfolioPosition{Quantity: 2, Value: 2, Open: 1, Close: 1, High: 1, Low: 1, AsOfDate: d1})
	sum.Positions = pl
	sum.CalcDailyTotals()

	assert.Equal(t, 2.0, sum.CurrentValue) //Does not get updated
	assert.Equal(t, 2.0, sum.CurrentOpen)
	assert.Equal(t, 2.0, sum.CurrentClose)
	assert.Equal(t, 2.0, sum.CurrentHigh)
	assert.Equal(t, 2.0, sum.CurrentLow)

	//Test update date (Should be the oldest date, which is d2)
	pl = append(pl, PortfolioPosition{Quantity: 2, Value: 2, Open: 1, Close: 1, High: 1, Low: 1, AsOfDate: d2})
	pl = append(pl, PortfolioPosition{Quantity: 2, Value: 2, Open: 1, Close: 1, High: 1, Low: 1, AsOfDate: d3})
	sum.Positions = pl
	sum.CalcDailyTotals()
	assert.Equal(t, d2, sum.AsOfDate)

	fmlogger.Exit(method)
}

func TestLoadPositions(t *testing.T) {
	method := "UserStockPortfolioSummary_test.TestLoadPositions"
	fmlogger.Enter(method)

	var sum UserStockPortfolioSummary
	var pl []PortfolioPosition
	d1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	pl = append(pl, PortfolioPosition{Quantity: 2, Value: 2, Open: 1, Close: 1, High: 1, Low: 1, AsOfDate: d1})
	sum.LoadPositions(pl)

	//Assert Totals were calculated
	assert.NotEqual(t, 0, sum.CurrentValue)
	assert.NotEqual(t, 0, sum.CurrentHigh)
	assert.NotEqual(t, 0, sum.CurrentLow)
	assert.NotEqual(t, 0, sum.CurrentOpen)
	assert.NotEqual(t, 0, sum.CurrentClose)
	assert.False(t, sum.AsOfDate.IsZero())

	//Assert positions were loaded
	assert.Equal(t, 1, len(sum.Positions))

	fmlogger.Enter(method)
}
