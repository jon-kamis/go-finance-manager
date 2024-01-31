package fmhandler

import (
	"encoding/json"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/test"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetUserStockPortfolioSummary_200(t *testing.T) {
	method := "summary_handler_test.TestGetUserStockPortfolioSummary_200"
	fmlogger.Enter(method)

	token := test.GetUserJWT(t)
	var resp models.UserStockPortfolioSummary

	setupStockTestData()

	writer := MakeRequest(http.MethodGet, "/users/2/stock-portfolio", nil, true, token)
	assert.Equal(t, http.StatusOK, writer.Code)

	//Read response values
	err := json.Unmarshal(writer.Body.Bytes(), &resp)
	assert.Nil(t, err)

	//Validate response values
	assert.Equal(t, 1, len(resp.Positions))
	assert.Equal(t, 2.0, resp.CurrentValue)

	teardownStocktestData()

	fmlogger.Exit(method)
}

func TestGetUserStockPortfolioSummary_403(t *testing.T) {
	method := "summary_handler_test.TestGetUserStockPortfolioSummary_403"
	fmlogger.Enter(method)

	token := test.GetUserJWT(t)

	writer := MakeRequest(http.MethodGet, "/users/1/stock-portfolio", nil, true, token)
	assert.Equal(t, http.StatusForbidden, writer.Code)

	fmlogger.Exit(method)
}

func setupStockTestData() {

	s1 := models.Stock{
		ID:           23,
		Ticker:       "AAPL",
		High:         1,
		Low:          1,
		Open:         1,
		Close:        1,
		Date:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CreateDt:     time.Now(),
		LastUpdateDt: time.Now(),
	}

	us1 := models.UserStock{
		ID:           23,
		UserId:       2,
		Ticker:       "AAPL",
		Quantity:     2,
		CreateDt:     time.Now(),
		LastUpdateDt: time.Now(),
	}

	p.GormDB.Create(&s1)
	fmh.DB.InsertUserStock(us1)

}

func teardownStocktestData() {

	s1 := models.Stock{
		ID: 23,
	}

	us1 := models.UserStock{
		ID: 23,
	}

	p.GormDB.Delete(s1)
	p.GormDB.Delete(us1)

}
