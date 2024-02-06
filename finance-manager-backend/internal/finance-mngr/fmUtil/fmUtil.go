// Package fmUtil contains utility methods used by the Finance Manager application
package fmUtil

import (
	"time"

	"github.com/jon-kamis/klogger"
)

// Function GetMonthBeginDate takes a date parameter and returns the first instant of the month
// that date is in
func GetMonthBeginDate(date time.Time) time.Time {
	method := "fmUtil.getMonthBeginDate"
	klogger.Enter(method)

	newDate := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)

	klogger.Exit(method)
	return newDate
}

// Function GetMonthEndDate takes a date parameter and returns the last instant of the month
// that date is in
func GetMonthEndDate(date time.Time) time.Time {
	method := "fmUtil.GetMonthEndDate"
	klogger.Enter(method)

	newDate := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	newDate = newDate.Add(-1 * time.Nanosecond)

	klogger.Exit(method)
	return newDate
}

// Function GetStartOfDay truncates a supplied time to the beginning of the day
func GetStartOfDay(date time.Time) time.Time {
	method := "fmUtil.GetStartOfDay"
	klogger.Enter(method)

	truncTime := 24 * time.Hour
	newDate := date.Truncate(truncTime)

	klogger.Exit(method)
	return newDate
}
