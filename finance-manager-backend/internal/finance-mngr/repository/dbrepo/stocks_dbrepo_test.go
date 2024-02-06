package dbrepo

import (
	"finance-manager-backend/internal/finance-mngr/models"
	"testing"

	"github.com/jon-kamis/klogger"
	"github.com/stretchr/testify/assert"
)

func TestInsertStock(t *testing.T) {
	method := "user_stocks_dbrepo_test.TestInsertStock"
	klogger.Enter(method)

	s := models.Stock{
		Ticker: "TEST1",
		High:   1,
		Low:    1,
		Open:   1,
		Close:  1,
	}

	id, err := d.InsertStock(s)
	assert.Nil(t, err)
	assert.Greater(t, id, 0)

	var sDb models.Stock
	err = p.GormDB.First(&sDb, id).Error

	assert.Nil(t, err)
	assert.Equal(t, id, sDb.ID)
	assert.Equal(t, s.Ticker, sDb.Ticker)
	assert.Equal(t, s.High, sDb.High)
	assert.Equal(t, s.Low, sDb.Low)

	//Clenaup
	p.GormDB.Delete(sDb)

	klogger.Exit(method)
}

func TestGetStockByTicker(t *testing.T) {
	method := "user_stocks_dbrepo_test.TestGetStockByTicker"
	klogger.Enter(method)

	setupStocks()

	sDb, err := d.GetStockByTicker("TEST1")
	assert.Nil(t, err)

	assert.Equal(t, "TEST1", sDb.Ticker)
	assert.NotEqual(t, 0, sDb.ID)

	//Test with something that does not exist
	sDb, err = d.GetStockByTicker("DOESNOTEXIST")
	assert.Nil(t, err)
	assert.Equal(t, 0, sDb.ID)

	tearDownStocks()
	klogger.Exit(method)
}

func setupStocks() {
	s := models.Stock{
		ID:     67,
		Ticker: "TEST1",
		High:   1,
		Low:    1,
		Open:   1,
		Close:  1,
	}
	p.GormDB.Create(&s)
}

func tearDownStocks() {

	s := models.Stock{
		ID: 67,
	}

	p.GormDB.Delete(s)
}
