package fmservice

import (
	"database/sql"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"fmt"
	"testing"
	"time"

	"github.com/jon-kamis/klogger"
	"github.com/stretchr/testify/assert"
)

func TestGetUserPortfolioBalanceHistory(t *testing.T) {
	method := "fm_stockservice.TestGetUserPortfolioBalanceHistory"
	klogger.Enter(method)

	d := time.Now() //This method pulls records based on current time, so we will use current time to set up data
	d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)

	//Test before data is entered
	hist, err := fms.GetUserPortfolioBalanceHistory(1, 5)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(hist))

	//Test with invalid userId
	hist, err = fms.GetUserPortfolioBalanceHistory(0, 5)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(hist))

	//Test with invalid date ranges
	hist, err = fms.GetUserPortfolioBalanceHistory(1, -1)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(hist))
	hist, err = fms.GetUserPortfolioBalanceHistory(1, 380)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(hist))

	//UserStockData to test with
	us1 := models.UserStock{
		ID:           22,
		UserId:       1,
		Type:         constants.UserStockTypeOwn,
		Ticker:       "AAPL",
		Quantity:     2,
		EffectiveDt:  d.Add(-5 * 24 * time.Hour),
		ExpirationDt: sql.NullTime{Time: d.Add(-1 * 24 * time.Hour)},
	}
	us2 := models.UserStock{
		ID:          23,
		UserId:      1,
		Type:        constants.UserStockTypeOwn,
		Ticker:      "MSFT",
		Quantity:    2,
		EffectiveDt: d.Add(-3 * 24 * time.Hour),
	}
	us3 := models.UserStock{
		ID:           24,
		UserId:       1,
		Type:         constants.UserStockTypeOwn,
		Ticker:       "MSFT",
		Quantity:     1,
		EffectiveDt:  d.Add(-4 * 24 * time.Hour),
		ExpirationDt: sql.NullTime{Time: d.Add(-3 * 24 * time.Hour).Add(-1 * time.Millisecond)},
	}

	fms.DB.InsertUserStock(us1)
	fms.DB.InsertUserStock(us2)
	fms.DB.InsertUserStock(us3)

	//StockData to test with
	sd := d.Add(-10 * 24 * time.Hour)

	for sd.Compare(d) <= 0 {
		s1 := models.StockData{
			Ticker: "MSFT",
			Close:  1,
			Date:   sd,
		}
		s2 := models.StockData{
			Ticker: "SNAP",
			Close:  1,
			Date:   sd,
		}
		s3 := models.StockData{
			Ticker: "AAPL",
			Close:  1,
			Date:   sd,
		}
		p.GormDB.Create(&s1)
		p.GormDB.Create(&s2)
		p.GormDB.Create(&s3)
		sd = sd.Add(24 * time.Hour)
	}

	hist, err = fms.GetUserPortfolioBalanceHistory(1, 5)

	dvm := make(map[time.Time]float64)

	for _, h := range hist {
		fmt.Printf("Adding date %v\n", h.Date)
		dvm[h.Date] = h.Close
	}

	assert.Nil(t, err)
	assert.Equal(t, 5, len(hist))
	assert.Equal(t, 2.0, dvm[d])                      //AAPL is expired, MSFT is worth 1 and quantity is x2
	assert.Equal(t, 4.0, dvm[d.Add(-24*time.Hour)])   //Both AAPL and MSFT are unexpired, both worth 1 with quantity x2
	assert.Equal(t, 3.0, dvm[d.Add(-24*4*time.Hour)]) //AAPL has quanity x2, MSFT has quantity x1, both have value of 1. Expect 3

	//Cleanup
	p.GormDB.Exec("DELETE FROM user_stocks")
	p.GormDB.Exec("DELETE FROM stock_data")

	klogger.Exit(method)
}
