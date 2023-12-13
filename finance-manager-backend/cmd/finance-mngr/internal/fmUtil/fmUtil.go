package fmUtil

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"time"
)

func GetMonthBeginDate(date time.Time) time.Time {
	method := "fmUtil.getMonthBeginDate"
	fmlogger.Enter(method)

	fmlogger.Exit(method)
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
}
