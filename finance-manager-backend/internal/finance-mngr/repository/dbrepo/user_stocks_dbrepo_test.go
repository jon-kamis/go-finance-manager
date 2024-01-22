package dbrepo

import (
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInsertUserStock(t *testing.T) {
	method := "user_stocks_dbrepo_test.TestInsertUserStock"
	fmlogger.Enter(method)

	s := models.UserStock{
		UserId:   1,
		Ticker:   "TEST1",
		Quantity: 2,
		EffectiveDt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	id, err := d.InsertUserStock(s)
	assert.Nil(t, err)
	assert.Greater(t, id, 0)

	var sDb models.UserStock
	err = p.GormDB.First(&sDb, id).Error

	assert.Nil(t, err)
	assert.True(t, sDb.ExpirationDt.Time.IsZero())
	assert.Equal(t, id, sDb.ID)
	assert.Equal(t, s.UserId, sDb.UserId)
	assert.Equal(t, s.Ticker, sDb.Ticker)
	assert.Equal(t, s.Quantity, sDb.Quantity)

	//Delete Record
	p.GormDB.Delete(sDb)

	var sDb2 models.UserStock
	s.ExpirationDt.Time = time.Now()
	id, err = d.InsertUserStock(s)
	assert.Nil(t, err)
	assert.Greater(t, id, 0)

	err = p.GormDB.Where("id=?",id).Find(&sDb2).Error

	assert.Nil(t, err)
	assert.False(t, sDb2.ExpirationDt.Time.IsZero())
	assert.Equal(t, id, sDb2.ID)

	p.GormDB.Delete(sDb2)

	fmlogger.Exit(method)
}

func TestGetAllUserStocks(t *testing.T) {
	method := "user_stocks_dbrepo_test.TestGetAllUserStocks"
	fmlogger.Enter(method)

	s1 := models.UserStock{
		UserId:      1,
		Ticker:      "TEST1",
		Quantity:    2,
		EffectiveDt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	s2 := models.UserStock{
		UserId:      1,
		Ticker:      "TEST2",
		Quantity:    2,
		EffectiveDt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	s3 := models.UserStock{
		UserId:   2,
		Ticker:   "TEST2",
		Quantity: 2,
		EffectiveDt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	//We need to insert using db method as gorm will incorrectly save expiration date as a zero date
	id1, _ := d.InsertUserStock(s1)
	id2, _ := d.InsertUserStock(s2)
	id3, _ := d.InsertUserStock(s3)

	//Get all user stocks with no search
	stocks, err := d.GetAllUserStocks(1, "", time.Now())
	assert.Nil(t, err)
	assert.NotNil(t, stocks)
	assert.Equal(t, 2, len(stocks))

	//Get all user stocks that have number 2 in ticker. Expect 1 result
	stocks, err = d.GetAllUserStocks(1, "2", time.Now())
	assert.Nil(t, err)
	assert.NotNil(t, stocks)
	assert.Equal(t, 1, len(stocks))

	//Get all user stocks for user with no data
	stocks, err = d.GetAllUserStocks(3, "", time.Now())
	assert.Nil(t, err)
	assert.NotNil(t, stocks)
	assert.Equal(t, 0, len(stocks))

	//Clean DB
	s1.ID = id1
	s2.ID = id2
	s3.ID = id3
	p.GormDB.Delete(s1)
	p.GormDB.Delete(s2)
	p.GormDB.Delete(s3)

	fmlogger.Exit(method)
}
