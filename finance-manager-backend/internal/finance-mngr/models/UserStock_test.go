package models

import (
	"testing"

	"github.com/jon-kamis/klogger"
	"github.com/stretchr/testify/assert"
)

func TestValidateCanSaveUserStock(t *testing.T) {
	method := "UserStock_test.TestValidateCanSaveUserStock"
	klogger.Enter(method)

	u := UserStock{
		UserId:   1,
		Quantity: 1,
		Ticker:   "AAPL",
	}

	err := u.ValidateCanSaveUserStock()
	assert.Nil(t, err)

	ut := u
	ut.UserId = 0
	err = ut.ValidateCanSaveUserStock()
	assert.NotNil(t, err)

	ut = u
	ut.Quantity = -1
	err = ut.ValidateCanSaveUserStock()
	assert.NotNil(t, err)

	ut = u
	ut.Ticker = ""
	err = ut.ValidateCanSaveUserStock()
	assert.NotNil(t, err)

	klogger.Exit(method)
}
