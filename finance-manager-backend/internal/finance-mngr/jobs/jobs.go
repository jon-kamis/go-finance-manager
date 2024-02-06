package jobs

import (
	"finance-manager-backend/internal/finance-mngr/application"
	"finance-manager-backend/internal/finance-mngr/constants"
	"time"

	"github.com/jon-kamis/klogger"
)

func ScheduledMinuteJobs(tick *time.Ticker, app application.Application) {
	method := "jobs.scheduleJobs"
	klogger.Info(method, "started running in asynchronous thread")

	updateStocks(time.Now(), app)
	for t := range tick.C {
		updateStocks(t, app)
	}
}

func updateStocks(t time.Time, app application.Application) {
	method := "jobs.updateStocks"
	klogger.Enter(method)

	if !app.ExternalService.GetIsStocksEnabled() {
		klogger.Info(method, "stocks are not enabled")
		klogger.Exit(method)
		return
	}

	s, err := app.DB.GetOldestStock()

	if err != nil {
		klogger.Error(method, constants.UnexpectedSQLError, err)
		klogger.Warn(method, "completed execution unsuccessfully")
		return
	}

	//only run if the date of s is at least 1 days old
	//tz, err := time.LoadLocation("EST")

	if err != nil {
		klogger.Error(method, constants.UnexpectedSQLError, err)
		klogger.Warn(method, "completed execution unsuccessfully")
		return
	}

	//We should check against the latest weekday

	compareDt := time.Now()
	compareDt = time.Date(compareDt.Year(), compareDt.Month(), compareDt.Day(), 0, 0, 0, 0, time.UTC)

	if compareDt.Weekday() == time.Monday {
		compareDt = compareDt.Add(-3 * 24 * time.Hour)
	} else if compareDt.Weekday() == time.Sunday {
		compareDt = compareDt.Add(-2 * 24 * time.Hour)
	} else {
		compareDt = compareDt.Add(-1 * 24 * time.Hour)
	}

	//get last date for stockData entries
	sd, err := app.DB.GetLatestStockDataByTicker(s.Ticker)

	if err != nil {
		klogger.Error(method, constants.UnexpectedSQLError, err)
		klogger.Warn(method, "completed execution unsuccessfully")
		return
	}

	klogger.Debug(method, "Checking if %v is before %v, or if stock data is not loaded for the given ticker", s.Date, compareDt)

	if s.Date.Before(compareDt) || sd.ID == 0 {

		var startDt time.Time

		//If sd is not loaded, then default to 1 year back
		if sd.ID == 0 {
			klogger.Debug(method, "stock data was not loaded. Loading stock data now")
			startDt = time.Now()
			startDt = time.Date(startDt.Year()-1, startDt.Month(), startDt.Day(), 0, 0, 0, 0, time.UTC)
		} else {
			startDt = sd.Date.Add(24 * time.Hour)
		}

		klogger.Debug(method, "fetching updates for stock "+s.Ticker)
		sn, err := app.ExternalService.FetchStockWithTickerForDateRange(s.Ticker, startDt, compareDt)

		if err != nil {
			klogger.Error(method, constants.UnexpectedSQLError, err)
			klogger.Warn(method, "completed execution unsuccessfully")
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
			klogger.Error(method, constants.UnexpectedSQLError, err)
			klogger.Warn(method, "completed execution unsuccessfully")
			return
		}

		//Save Stock Data
		err = app.DB.InsertStockData(sn)

		if err != nil {
			klogger.Error(method, constants.UnexpectedSQLError, err)
			klogger.Warn(method, "completed execution unsuccessfully")
			return
		}

	} else {
		klogger.Info(method, "oldest stock is up to date")
	}

	klogger.Exit(method)
}
