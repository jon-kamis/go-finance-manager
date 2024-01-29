package fmstockservice

import (
	"database/sql"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/internal/finance-mngr/repository/dbrepo"
	"finance-manager-backend/test"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var p test.DockerTestPlatform
var fss FmStockService

func TestMain(m *testing.M) {
	method := "validation_test.TestMain"
	fmlogger.Enter(method)

	//Setup testing platform using docker
	p = test.Setup(m)
	db := &dbrepo.PostgresDBRepo{DB: p.DB}

	fss = FmStockService{
		StocksEnabled: false,
		DB:            db,
	}

	//Execute Code
	code := m.Run()

	//Tear down testing platform
	test.TearDown(p)

	fmlogger.Exit(method)
	os.Exit(code)
}
func TestLoadApiKeyFromFile(t *testing.T) {
	method := "FinanceManagerHandler_test.TestLoadApiKeyFromFile"
	fmlogger.Enter(method)

	//Write File to read from
	pwd, _ := os.Getwd()
	fileName := "TestLoadApiKeyFromFile.keytest"
	content := "test content"

	err := os.WriteFile(pwd+fileName, []byte(content), 0666)

	if err != nil {
		t.Errorf("failed to persist file prior to test")
	}

	fss := FmStockService{
		StocksApiKeyFileName: fileName,
		StocksEnabled:        false,
	}

	//Run Test
	err = fss.LoadApiKeyFromFile()
	assert.Nil(t, err)
	assert.True(t, fss.StocksEnabled)
	assert.Equal(t, content, fss.PolygonApiKey)

	//Run failing test
	fss.StocksApiKeyFileName = "someotherfile.keytest"
	err = fss.LoadApiKeyFromFile()
	assert.NotNil(t, err)

	//Clean up test file
	err = os.Remove(pwd + fileName)

	if err != nil {
		t.Errorf("failed to clean up test files after test completion")
	}

	fmlogger.Exit(method)
}

func TestUpdateAndPersistAPIKey(t *testing.T) {
	method := "fm_stockservice.TestLoadApiKeyFromFile"
	fmlogger.Enter(method)

	pwd, _ := os.Getwd()
	fileName := "TestUpdateAndPersistAPIKey.keytest"
	content := "test content"

	fss := FmStockService{
		StocksApiKeyFileName: fileName,
		StocksEnabled:        false,
	}

	//Run Test
	err := fss.UpdateAndPersistAPIKey(content)
	assert.Nil(t, err)
	assert.True(t, fss.StocksEnabled)
	assert.Equal(t, content, fss.PolygonApiKey)

	//Verify that test was successful
	_, err = os.ReadFile(pwd + fileName)
	assert.Nil(t, err)

	//Clean up the test file
	err = os.Remove(pwd + fileName)

	if err != nil {
		t.Errorf("failed to clean up test files after test completion")
	}

	fmlogger.Exit(method)
}

func TestGetUserPortfolioBalanceHistory(t *testing.T) {
	method := "fm_stockservice.TestGetUserPortfolioBalanceHistory"
	fmlogger.Enter(method)

	d := time.Now() //This method pulls records based on current time, so we will use current time to set up data
	d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)

	//Test before data is entered
	hist, err := fss.GetUserPortfolioBalanceHistory(1, 5)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(hist))

	//Test with invalid userId
	hist, err = fss.GetUserPortfolioBalanceHistory(0, 5)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(hist))

	//Test with invalid date ranges
	hist, err = fss.GetUserPortfolioBalanceHistory(1, -1)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(hist))
	hist, err = fss.GetUserPortfolioBalanceHistory(1, 380)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(hist))

	//UserStockData to test with
	us1 := models.UserStock{
		ID:           22,
		UserId:       1,
		Ticker:       "AAPL",
		Quantity:     2,
		EffectiveDt:  d.Add(-5 * 24 * time.Hour),
		ExpirationDt: sql.NullTime{Time: d.Add(-1 * 24 * time.Hour)},
	}
	us2 := models.UserStock{
		ID:          23,
		UserId:      1,
		Ticker:      "MSFT",
		Quantity:    2,
		EffectiveDt: d.Add(-3 * 24 * time.Hour),
	}
	us3 := models.UserStock{
		ID:           24,
		UserId:       1,
		Ticker:       "MSFT",
		Quantity:     1,
		EffectiveDt:  d.Add(-4 * 24 * time.Hour),
		ExpirationDt: sql.NullTime{Time: d.Add(-3 * 24 * time.Hour).Add(-1 * time.Millisecond)},
	}

	fss.DB.InsertUserStock(us1)
	fss.DB.InsertUserStock(us2)
	fss.DB.InsertUserStock(us3)

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

	hist, err = fss.GetUserPortfolioBalanceHistory(1, 5)

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

	fmlogger.Exit(method)
}
