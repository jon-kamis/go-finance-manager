package main

import (
	"finance-manager-backend/cmd/finance-mngr/internal/application"
	"finance-manager-backend/cmd/finance-mngr/internal/constants"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"fmt"
	"time"
)

func scheduledMinuteJobs(tick *time.Ticker, app application.Application) {
	method := "jobs.scheduleJobs"
	fmlogger.Enter(method)

	updateStocks(time.Now(), app)
	for t := range tick.C {
		updateStocks(t, app)
	}

	fmlogger.Exit(method)
}

func updateStocks(t time.Time, app application.Application) {
	method := "jobs.updateStocks"
	fmt.Printf("[%s] began execution at %v\n", method, t)

	if !app.StocksService.GetIsStocksEnabled() {
		fmlogger.Info(method, "stocks are not enabled")
		fmt.Printf("[%s] completed execution at %v\n", method, time.Now())
		return
	}

	s, err := app.DB.GetOldestStock()

	if err != nil {
		fmlogger.Error(method, constants.UnexpectedSQLError, err)
		fmt.Printf("[%s] completed execution unsuccessfully at %v\n", method, time.Now())
		return
	}

	//only run if the date of s is at least 1 days old
	limit := time.Now().Add(-29 * time.Hour)

	if s.Date.Before(limit) {
		fmlogger.Info(method, "performing update")
		sn, err := app.StocksService.FetchStockWithTicker(s.Ticker)
		
		if err != nil {
			fmlogger.Error(method, constants.UnexpectedSQLError, err)
			fmt.Printf("[%s] completed execution unsuccessfully at %v\n", method, time.Now())
			return
		}

		s.Open = sn.Open
		s.Close = sn.Close
		s.High = sn.High
		s.Low = sn.Low
		s.Date = sn.Date

		//Save updated values
		err = app.DB.UpdateStock(s)

		if err != nil {
			fmlogger.Error(method, constants.UnexpectedSQLError, err)
			fmt.Printf("[%s] completed execution unsuccessfully at %v\n", method, time.Now())
			return
		}

	} else {
		fmlogger.Info(method, "oldest stock is up to date")
	}

	fmt.Printf("[%s] completed execution at %v\n", method, time.Now())
}
