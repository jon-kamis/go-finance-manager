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
