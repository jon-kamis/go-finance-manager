package restmodels

import (
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/enums/stockoperation"
	"testing"
	"time"

	"github.com/jon-kamis/klogger"
	"github.com/stretchr/testify/assert"
)

func TestIsValidRequest(t *testing.T) {
	method := "ModifyStockRequest_test.TestIsValidRequest"
	klogger.Enter(method)

	r := ModifyStockRequest{
		Ticker:    "AAPL",
		Amount:    5,
		Operation: stockoperation.Add,
		Date:      time.Now(),
	}

	var r1 ModifyStockRequest
	var v bool
	var m string

	v, m = r.IsValidRequest()
	assert.True(t, v)
	assert.Equal(t, "", m)

	r1 = r
	r1.Ticker = ""
	v, m = r1.IsValidRequest()
	assert.False(t, v)
	assert.Equal(t, constants.StockOperationTickerRequiredError, m)

	r1 = r
	r1.Amount = 0
	v, m = r1.IsValidRequest()
	assert.False(t, v)
	assert.Equal(t, constants.StockOperationInvalidAmountError, m)

	r1 = r
	r1.Operation = "Invalid"
	v, m = r1.IsValidRequest()
	assert.False(t, v)
	assert.Equal(t, constants.StockOperationInvalidOperationError, m)

	r1 = r
	var tz time.Time
	r1.Date = tz
	v, m = r1.IsValidRequest()
	assert.False(t, v)
	assert.Equal(t, constants.StockOperationInvalidDateError, m)

	r1 = r
	r1.Date = time.Now().Add(24 * time.Hour)
	v, m = r1.IsValidRequest()
	assert.False(t, v)
	assert.Equal(t, constants.StockOperationInvalidDateError, m)

	klogger.Exit(method)
}
