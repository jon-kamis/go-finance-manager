package dbrepo

import (
	"database/sql"
	"finance-manager-backend/internal/finance-mngr/constants"
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
		UserId:      1,
		Ticker:      "TEST1",
		Quantity:    2,
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

	err = p.GormDB.Where("id=?", id).Find(&sDb2).Error

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
		Type:        constants.UserStockTypeOwn,
		Ticker:      "TEST1",
		Quantity:    2,
		EffectiveDt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	s2 := models.UserStock{
		UserId:      1,
		Type:        constants.UserStockTypeOwn,
		Ticker:      "TEST2",
		Quantity:    2,
		EffectiveDt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	s3 := models.UserStock{
		UserId:      2,
		Type:        constants.UserStockTypeOwn,
		Ticker:      "TEST2",
		Quantity:    2,
		EffectiveDt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	s4 := models.UserStock{
		UserId:      1,
		Type:        constants.UserStockTypeWatch,
		Ticker:      "TEST2",
		Quantity:    2,
		EffectiveDt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	//We need to insert using db method as gorm will incorrectly save expiration date as a zero date
	id1, _ := d.InsertUserStock(s1)
	id2, _ := d.InsertUserStock(s2)
	id3, _ := d.InsertUserStock(s3)
	id4, _ := d.InsertUserStock(s4)

	//Get all user stocks with no search
	stocks, err := d.GetAllUserStocks(1, "", "", time.Now())
	assert.Nil(t, err)
	assert.NotNil(t, stocks)
	assert.Equal(t, 2, len(stocks))

	//Get all user stocks that have number 2 in ticker. Expect 1 result
	stocks, err = d.GetAllUserStocks(1, "", "2", time.Now())
	assert.Nil(t, err)
	assert.NotNil(t, stocks)
	assert.Equal(t, 1, len(stocks))

	//Get user stock watchlist
	stocks, err = d.GetAllUserStocks(1, constants.UserStockTypeWatch, "", time.Now())
	assert.Nil(t, err)
	assert.NotNil(t, stocks)
	assert.Equal(t, 1, len(stocks))

	//Get all user stocks for user with no data
	stocks, err = d.GetAllUserStocks(3, "", "", time.Now())
	assert.Nil(t, err)
	assert.NotNil(t, stocks)
	assert.Equal(t, 0, len(stocks))

	//Clean DB
	s1.ID = id1
	s2.ID = id2
	s3.ID = id3
	s4.ID = id4
	p.GormDB.Delete(s1)
	p.GormDB.Delete(s2)
	p.GormDB.Delete(s3)
	p.GormDB.Delete(s4)

	fmlogger.Exit(method)
}

func TestGetAllUserStocksByDateRange(t *testing.T) {
	method := "user_stocks_dbrepo_test.TestGetAllUserStocksByDateRange"
	fmlogger.Enter(method)

	s1 := models.UserStock{
		ID:           17,
		UserId:       1,
		Type:         constants.UserStockTypeOwn,
		Ticker:       "AAPL",
		Quantity:     2,
		EffectiveDt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
		ExpirationDt: sql.NullTime{Time: time.Date(2024, 1, 3, 23, 59, 59, 999, time.Local), Valid: true},
	}

	s2 := models.UserStock{
		ID:          18,
		UserId:      1,
		Type:        constants.UserStockTypeOwn,
		Ticker:      "MSFT",
		Quantity:    2,
		EffectiveDt: time.Date(2024, 1, 4, 0, 0, 0, 0, time.Local),
		ExpirationDt: sql.NullTime{Time: time.Date(2024, 1, 7, 23, 59, 59, 999, time.Local), Valid: true},
	}

	s3 := models.UserStock{
		ID:          19,
		UserId:      1,
		Type:        constants.UserStockTypeOwn,
		Ticker:      "SNAP",
		Quantity:    2,
		EffectiveDt: time.Date(2024, 1, 8, 0, 0, 0, 0, time.Local),
	}

	//Insert Test data
	p.GormDB.Create(&s1)
	p.GormDB.Create(&s2)
	p.GormDB.Create(&s3)

	var usl []*models.UserStock
	var err error

	usl, err = d.GetAllUserStocksByDateRange(1, "", time.Date(2023, 12, 30, 0, 0, 0, 0, time.Local), time.Date(2023, 12, 31, 0, 0, 0, 0, time.Local))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(usl))

	usl, err = d.GetAllUserStocksByDateRange(1, "", time.Date(2023, 12, 31, 0, 0, 0, 0, time.Local), time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(usl))
	assert.Equal(t, 17, usl[0].ID)

	usl, err = d.GetAllUserStocksByDateRange(1, "", time.Date(2023, 12, 31, 0, 0, 0, 0, time.Local), time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local))
	assert.Nil(t, err)
	assert.Equal(t, 3, len(usl))

	usl, err = d.GetAllUserStocksByDateRange(1, "", time.Date(2024, 12, 31, 0, 0, 0, 0, time.Local), time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(usl))
	assert.Equal(t, 19, usl[0].ID)

	//Clean Test Data
	p.GormDB.Delete(&s1)
	p.GormDB.Delete(&s2)
	p.GormDB.Delete(&s3)

	fmlogger.Exit(method)
}

func TestUpdateUserStock(t *testing.T) {
	method := "user_stocks_dbrepo_test.TestUpdateUserStock"
	fmlogger.Enter(method)

	s1 := models.UserStock{
		ID:          27,
		UserId:      1,
		Type:        constants.UserStockTypeOwn,
		Ticker:      "TEST1",
		Quantity:    2,
		EffectiveDt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	s2 := models.UserStock{
		ID:          27,
		UserId:      1,
		Type:        constants.UserStockTypeOwn,
		Ticker:      "TEST2",
		Quantity:    3,
		EffectiveDt: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	//Insert Record to start
	p.GormDB.Create(&s1)

	//Update record
	err := d.UpdateUserStock(s2)
	assert.Nil(t, err)

	var sdb models.UserStock
	err = p.GormDB.Exec("SELECT * FROM user_stocks WHERE id = 27").Find(&sdb).Error

	assert.Nil(t, err)
	assert.Equal(t, 3.0, sdb.Quantity)
	assert.Equal(t, s1.Ticker, sdb.Ticker) //Ticker does not get updated

	//Cleanup Record
	p.GormDB.Delete(&s1)

	fmlogger.Exit(method)
}

func TestGetUserStockByUserIdTickerAndDate(t *testing.T) {
	method := "user_stocks_dbrepo_test.TestGetUserStockByUserIdTickerAndDate"
	fmlogger.Enter(method)

	s1 := models.UserStock{
		ID:           37,
		UserId:       1,
		Type:         constants.UserStockTypeOwn,
		Ticker:       "AAPL",
		Quantity:     2,
		EffectiveDt:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
		ExpirationDt: sql.NullTime{Time: time.Date(2024, 1, 1, 23, 59, 59, 999, time.Local), Valid: true},
	}

	s2 := models.UserStock{
		ID:          38,
		UserId:      1,
		Type:        constants.UserStockTypeOwn,
		Ticker:      "AAPL",
		Quantity:    2,
		EffectiveDt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.Local),
	}

	p.GormDB.Create(&s1)
	p.GormDB.Create(&s2)

	//Between eff and exp date for a record
	sdb, err := d.GetUserStockByUserIdTickerAndDate(1, "AAPL", time.Date(2023, 12, 31, 2, 1, 43, 234, time.Local))
	assert.Nil(t, err)
	assert.Equal(t, 37, sdb.ID)

	//After last record with no exp date
	sdb, err = d.GetUserStockByUserIdTickerAndDate(1, "AAPL", time.Date(2025, 12, 31, 2, 1, 43, 234, time.Local))
	assert.Nil(t, err)
	assert.Equal(t, 38, sdb.ID)

	//Before first record
	sdb, err = d.GetUserStockByUserIdTickerAndDate(1, "AAPL", time.Date(2022, 12, 31, 2, 1, 43, 234, time.Local))
	assert.Nil(t, err)
	assert.Equal(t, 0, sdb.ID)

	//Clean data
	p.GormDB.Delete(&s1)
	p.GormDB.Delete(&s2)

	fmlogger.Exit(method)
}
