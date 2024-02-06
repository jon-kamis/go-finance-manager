package fmservice

import (
	"database/sql"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/enums/stockoperation"
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/internal/finance-mngr/models/restmodels"
	"testing"
	"time"

	"github.com/jon-kamis/klogger"
	"github.com/stretchr/testify/assert"
)

func TestLoadPriorUserStockForTransaction(t *testing.T) {
	method := "user_stock_service_test.TestLoadPriorUserStockForTransaction"
	klogger.Enter(method)

	createTestLoadPriorUserStockForTransactionData()

	d := time.Date(2024, 1, 10, 0, 0, 0, 0, time.Local)

	r := restmodels.ModifyStockRequest{
		Ticker:    "AAPL",
		Amount:    5,
		Operation: stockoperation.Add,
		Date:      d,
	}

	var p models.UserStock
	u := models.UserStock{UserId: 1}

	//Test adding to existing user stock, between transactions
	err := fms.LoadPriorUserStockForTransaction(r, &p, &u)
	assert.Nil(t, err)

	assert.Equal(t, d.Add(-1*time.Millisecond), p.ExpirationDt.Time)
	assert.Equal(t, time.Date(2024, 1, 13, 23, 59, 59, 0, time.Local), u.ExpirationDt.Time)
	assert.Equal(t, 15.0, u.Quantity)

	//Test removing from existing user stock, between transactions
	r.Operation = stockoperation.Remove
	err = fms.LoadPriorUserStockForTransaction(r, &p, &u)
	assert.Nil(t, err)

	assert.Equal(t, d.Add(-1*time.Millisecond), p.ExpirationDt.Time)
	assert.Equal(t, time.Date(2024, 1, 13, 23, 59, 59, 0, time.Local), u.ExpirationDt.Time)
	assert.Equal(t, 5.0, u.Quantity)

	//Test removing from existing user stock, between transactions where quantity would fall below 0
	r.Operation = stockoperation.Remove
	r.Amount = 15.0
	err = fms.LoadPriorUserStockForTransaction(r, &p, &u)
	assert.NotNil(t, err)

	//Test removing from nonexisting user stock
	r.Operation = stockoperation.Remove
	r.Date = time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)
	r.Amount = 15.0
	err = fms.LoadPriorUserStockForTransaction(r, &p, &u)
	assert.NotNil(t, err)

	//Test adding to nonexisting user stock
	r.Operation = stockoperation.Add
	r.Date = time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)
	r.Amount = 15.0
	err = fms.LoadPriorUserStockForTransaction(r, &p, &u)
	assert.Nil(t, err)
	assert.Equal(t, 15.0, u.Quantity)
	assert.Equal(t, 0, p.ID)

	tearDownTestLoadPriorUserStockForTransactionData()

	klogger.Exit(method)
}

func createTestLoadPriorUserStockForTransactionData() {
	us1 := models.UserStock{
		ID:          1,
		UserId:      1,
		Ticker:      "AAPL",
		Quantity:    10,
		Type:        constants.UserStockTypeOwn,
		EffectiveDt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
		ExpirationDt: sql.NullTime{
			Time:  time.Date(2024, 1, 13, 23, 59, 59, 999, time.Local),
			Valid: true,
		},
	}

	us2 := models.UserStock{
		ID:          2,
		UserId:      1,
		Ticker:      "AAPL",
		Quantity:    10,
		Type:        constants.UserStockTypeOwn,
		EffectiveDt: time.Date(2024, 1, 14, 0, 0, 0, 0, time.Local),
		ExpirationDt: sql.NullTime{
			Time:  time.Date(2024, 1, 31, 23, 59, 59, 999, time.Local),
			Valid: true,
		},
	}

	us3 := models.UserStock{
		ID:           3,
		UserId:       1,
		Ticker:       "AAPL",
		Quantity:     10,
		Type:         constants.UserStockTypeOwn,
		EffectiveDt:  time.Date(2024, 2, 1, 0, 0, 0, 0, time.Local),
		ExpirationDt: sql.NullTime{},
	}

	us4 := models.UserStock{
		ID:           4,
		UserId:       1,
		Ticker:       "MSFT",
		Quantity:     10,
		Type:         constants.UserStockTypeOwn,
		EffectiveDt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
		ExpirationDt: sql.NullTime{},
	}

	p.GormDB.Create(&us1)
	p.GormDB.Create(&us2)
	p.GormDB.Create(&us3)
	p.GormDB.Create(&us4)
}

func tearDownTestLoadPriorUserStockForTransactionData() {
	us1 := models.UserStock{
		ID: 1,
	}

	us2 := models.UserStock{
		ID: 2,
	}

	us3 := models.UserStock{
		ID: 3,
	}

	us4 := models.UserStock{
		ID: 4,
	}

	p.GormDB.Delete(&us1)
	p.GormDB.Delete(&us2)
	p.GormDB.Delete(&us3)
	p.GormDB.Delete(&us4)
}
