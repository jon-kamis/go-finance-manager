package fmhandler

import (
	"encoding/json"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/test"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetUserStocks_200(t *testing.T) {
	method := "stocks_handler_test.TestGetUserSTocks"
	fmlogger.Enter(method)

	token := test.GetUserJWTWithId(t, 3)
	var resp []models.Stock

	setupStockHandlerTestData()

	writer := MakeRequest(http.MethodGet, "/users/3/stocks", nil, true, token)
	assert.Equal(t, http.StatusOK, writer.Code)

	//Read response values
	err := json.Unmarshal(writer.Body.Bytes(), &resp)
	assert.Nil(t, err)

	//Validate response values
	assert.Equal(t, 1, len(resp))
	assert.Equal(t, "MSFT", resp[0].Ticker)

	teardownStockHandlertestData()

	fmlogger.Exit(method)
}

func TestGetUserStocks_403(t *testing.T) {
	method := "stocks_handler_test.TestGetUserSTocks"
	fmlogger.Enter(method)

	token := test.GetUserJWTWithId(t, 3)

	writer := MakeRequest(http.MethodGet, "/users/2/stocks", nil, true, token)
	assert.Equal(t, http.StatusForbidden, writer.Code)

	fmlogger.Exit(method)
}

func setupStockHandlerTestData() {

	s1 := models.Stock{
		ID:           24,
		Ticker:       "MSFT",
		High:         1,
		Low:          1,
		Open:         1,
		Close:        1,
		Date:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CreateDt:     time.Now(),
		LastUpdateDt: time.Now(),
	}

	us1 := models.UserStock{
		ID:           24,
		UserId:       3,
		Ticker:       "MSFT",
		Quantity:     2,
		Type:         constants.UserStockTypeOwn,
		CreateDt:     time.Now(),
		LastUpdateDt: time.Now(),
		EffectiveDt:  time.Now().Add(-24 * time.Hour),
	}

	p.GormDB.Create(&s1)
	fmh.DB.InsertUserStock(us1)

}

func teardownStockHandlertestData() {

	s1 := models.Stock{
		ID: 24,
	}

	us1 := models.UserStock{
		ID: 24,
	}

	p.GormDB.Delete(s1)
	p.GormDB.Delete(us1)

}
