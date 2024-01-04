package fmUtil

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"testing"
	"time"
)

func TestGetMonthBeginDate(t *testing.T) {
	method := "fmUtil_test.TestMonthBeginDate"
	fmlogger.Enter(method)

	date := time.Date(2024, 2, 23, 7, 33, 32, 1, time.UTC)

	monthBeginDate := GetMonthBeginDate(date)

	if monthBeginDate.Month() != date.Month() || monthBeginDate.Year() != date.Year() || monthBeginDate.Day() != 1 || monthBeginDate.Hour() != 0 || monthBeginDate.Second() != 0 || monthBeginDate.Second() != 0 || monthBeginDate.Nanosecond() != 0 {
		fmlogger.Exit(method)
		t.Errorf("unexpected date time returned")
	}

	fmlogger.Exit(method)

}
