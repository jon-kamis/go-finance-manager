package fmUtil

import (
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMonthBeginDate(t *testing.T) {
	method := "fmUtil_test.TestMonthBeginDate"
	fmlogger.Enter(method)

	date := time.Date(2024, 2, 23, 7, 33, 32, 1, time.UTC)
	expectedDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	monthBeginDate := GetMonthBeginDate(date)

	assert.Equal(t, expectedDate, monthBeginDate)

	fmlogger.Exit(method)

}

func TestGetMonthEndDate(t *testing.T) {
	method := "fmUtil_test.TestMonthBeginDate"
	fmlogger.Enter(method)

	date := time.Date(2024, 2, 23, 7, 33, 32, 1, time.UTC)
	expectedDate := time.Date(2024, 2, 29, 23, 59, 59, 999999999, time.UTC)
	monthBeginDate := GetMonthEndDate(date)

	assert.Equal(t, expectedDate, monthBeginDate)

	fmlogger.Exit(method)

}

func TestGetStartOfDay(t *testing.T) {
	method := "fmUtil_test.TestGetStartOfDay"
	fmlogger.Enter(method)

	date := time.Date(2024, 2, 23, 7, 33, 32, 1, time.UTC)
	expectedTime := time.Date(2024, 2, 23, 0, 0, 0, 0, time.UTC)
	dayBeginTime := GetStartOfDay(date)

	assert.Equal(t, expectedTime, dayBeginTime)

	fmlogger.Exit(method)
}
