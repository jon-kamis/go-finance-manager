package jobs

import (
	"finance-manager-backend/internal/finance-mngr/application"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"fmt"
	"time"
)

func ScheduledMinuteJobs(tick *time.Ticker, app application.Application) {
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
	//tz, err := time.LoadLocation("EST")

	if err != nil {
		fmlogger.Error(method, constants.UnexpectedSQLError, err)
		fmt.Printf("[%s] completed execution unsuccessfully at %v\n", method, time.Now())
		return
	}

	yesterday := time.Now()
	yesterday = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
	yesterday = yesterday.Add(-24 * time.Hour)

	fmt.Printf("Checking if %v is before %v\n", s.Date, yesterday)

	if s.Date.Before(yesterday) {

		//get last date for stockData entries
		sd, err := app.DB.GetLatestStockDataByTicker(s.Ticker)

		if err != nil {
			fmlogger.Error(method, constants.UnexpectedSQLError, err)
			fmt.Printf("[%s] completed execution unsuccessfully at %v\n", method, time.Now())
			return
		}

		var startDt time.Time

		//If sd is not loaded, then default to 1 year back
		if sd.ID == 0 {
			startDt = time.Now()
			startDt = time.Date(startDt.Year()-1, startDt.Month(), startDt.Day(), 0, 0, 0, 0, time.UTC)
		} else {
			startDt = sd.Date
		}

		fmlogger.Info(method, "fetching updates for stock "+s.Ticker)
		sn, err := app.StocksService.FetchStockWithTickerForDateRange(s.Ticker, startDt, time.Now())

		if err != nil {
			fmlogger.Error(method, constants.UnexpectedSQLError, err)
			fmt.Printf("[%s] completed execution unsuccessfully at %v\n", method, time.Now())
			return
		}

		//Latest index should be most up to date entry
		i := len(sn) - 1

		s.Open = sn[i].Open
		s.Close = sn[i].Close
		s.High = sn[i].High
		s.Low = sn[i].Low
		s.Date = sn[i].Date

		//Save stock
		err = app.DB.UpdateStock(s)

		if err != nil {
			fmlogger.Error(method, constants.UnexpectedSQLError, err)
			fmt.Printf("[%s] completed execution unsuccessfully at %v\n", method, time.Now())
			return
		}

		//Save Stock Data
		err = app.DB.InsertStockData(sn)

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
