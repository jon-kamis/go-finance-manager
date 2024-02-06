package fmUtil

import (
	"finance-manager-backend/test/logtest"
	"os"
	"testing"
	"time"

	"github.com/jon-kamis/klogger"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	logtest.SetKloggerTestFileNameEnv()

	method := "auth_test.TestMain"
	klogger.Enter(method)

	//Execute Code
	code := m.Run()

	klogger.Exit(method)
	os.Exit(code)
}

func TestGetMonthBeginDate(t *testing.T) {
	method := "fmUtil_test.TestMonthBeginDate"
	klogger.Enter(method)

	date := time.Date(2024, 2, 23, 7, 33, 32, 1, time.UTC)
	expectedDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	monthBeginDate := GetMonthBeginDate(date)

	assert.Equal(t, expectedDate, monthBeginDate)

	klogger.Exit(method)

}

func TestGetMonthEndDate(t *testing.T) {
	method := "fmUtil_test.TestMonthBeginDate"
	klogger.Enter(method)

	date := time.Date(2024, 2, 23, 7, 33, 32, 1, time.UTC)
	expectedDate := time.Date(2024, 2, 29, 23, 59, 59, 999999999, time.UTC)
	monthBeginDate := GetMonthEndDate(date)

	assert.Equal(t, expectedDate, monthBeginDate)

	klogger.Exit(method)

}

func TestGetStartOfDay(t *testing.T) {
	method := "fmUtil_test.TestGetStartOfDay"
	klogger.Enter(method)

	date := time.Date(2024, 2, 23, 7, 33, 32, 1, time.UTC)
	expectedTime := time.Date(2024, 2, 23, 0, 0, 0, 0, time.UTC)
	dayBeginTime := GetStartOfDay(date)

	assert.Equal(t, expectedTime, dayBeginTime)

	klogger.Exit(method)
}
