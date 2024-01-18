package dbrepo

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertUserStock(t *testing.T) {
	method := "user_stocks_dbrepo_test.TestInsertUserStock"
	fmlogger.Enter(method)

	s := models.UserStock{
		UserId:   1,
		Ticker:   "TEST1",
		Quantity: 2,
	}

	id, err := d.InsertUserStock(s)
	assert.Nil(t, err)
	assert.Greater(t, id, 0)

	var sDb models.UserStock
	err = p.GormDB.First(&sDb, id).Error

	assert.Nil(t, err)
	assert.Equal(t, id, sDb.ID)
	assert.Equal(t, s.UserId, sDb.UserId)
	assert.Equal(t, s.Ticker, sDb.Ticker)
	assert.Equal(t, s.Quantity, sDb.Quantity)

	//Clenaup
	p.GormDB.Delete(sDb)

	fmlogger.Exit(method)
}

func TestGetAllUserStocks(t *testing.T) {
	method := "user_stocks_dbrepo_test.TestGetAllUserStocks"
	fmlogger.Enter(method)

	s1 := models.UserStock{
		UserId:   1,
		Ticker:   "TEST1",
		Quantity: 2,
	}

	s2 := models.UserStock{
		UserId:   1,
		Ticker:   "TEST2",
		Quantity: 2,
	}

	s3 := models.UserStock{
		UserId:   2,
		Ticker:   "TEST2",
		Quantity: 2,
	}

	p.GormDB.Create(&s1)
	p.GormDB.Create(&s2)
	p.GormDB.Create(&s3)

	//Get all user stocks with no search
	stocks, err := d.GetAllUserStocks(1, "")
	assert.Nil(t, err)
	assert.NotNil(t, stocks)
	assert.Equal(t, 2, len(stocks))

	//Get all user stocks that have number 2 in ticker. Expect 1 result
	stocks, err = d.GetAllUserStocks(1, "2")
	assert.Nil(t, err)
	assert.NotNil(t, stocks)
	assert.Equal(t, 1, len(stocks))

	//Get all user stocks for user with no data
	stocks, err = d.GetAllUserStocks(3, "")
	assert.Nil(t, err)
	assert.NotNil(t, stocks)
	assert.Equal(t, 0, len(stocks))

	//Clean DB
	p.GormDB.Delete(s1)
	p.GormDB.Delete(s2)
	p.GormDB.Delete(s3)

	fmlogger.Exit(method)
}
