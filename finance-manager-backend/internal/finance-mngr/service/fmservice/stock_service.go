package fmservice

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"sort"
	"time"

	"github.com/jon-kamis/klogger"
)

// Function GetUserPortfolioBalanceHistory fetches the Portfolio Balance History for a user for a given timeframe
// uId - The ID of the user to fetch history for
// d - The number of past days to pull history for. Maximum is 365
func (fms *FMService) GetUserPortfolioBalanceHistory(uId int, d int) ([]models.PortfolioBalanceHistory, error) {
	method := "fm_stockservice.GetUserPortfolioBalanceHistory"
	klogger.Enter(method)

	var hist []models.PortfolioBalanceHistory
	var sd time.Time
	var err error
	ed := time.Now()

	//Validate passed values
	if uId <= 0 {
		err = errors.New("uId is required")
		klogger.ExitError(method, err.Error())
		return hist, err
	}

	if d < 1 || d > 365 {
		err = errors.New("d must be between 1 and 365 inclusively")
		klogger.ExitError(method, err.Error())
		return hist, err
	}

	//First Load User Positions for date range
	sd = ed.Add(-1 * time.Duration(d) * 24 * time.Hour)

	usl, err := fms.DB.GetAllUserStocksByDateRange(uId, "", sd, ed)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return hist, err
	}

	//No stocks exist
	if len(usl) == 0 {
		klogger.Info(method, "User has no stocks for this timeframe")
		klogger.Exit(method)
		return hist, nil
	}

	//Next Loop through each user position and load stock data for that position. Add total value for each date
	histMap := make(map[time.Time]models.PortfolioBalanceHistory)

	for _, us := range usl {

		var d1 time.Time
		var d2 time.Time

		if us.EffectiveDt.Before(sd) {
			d1 = sd
		} else {
			d1 = us.EffectiveDt
		}

		if !us.ExpirationDt.Time.IsZero() && us.ExpirationDt.Time.Before(ed) {
			d2 = us.ExpirationDt.Time
		} else {
			d2 = ed
		}

		sl, err := fms.DB.GetStockDataByTickerAndDateRange(us.Ticker, d1, d2)

		if err != nil {
			klogger.ExitError(method, constants.UnexpectedSQLError, err)
			return hist, err
		}

		//Next, Loop through stock Data for this entry and add totals to each date in map
		for _, s := range sl {
			if histMap[s.Date].Date.IsZero() {

				//Initialize value
				hd := models.PortfolioBalanceHistory{
					Date:  s.Date,
					Close: us.Quantity * s.Close,
					Open:  us.Quantity * s.Open,
					High:  us.Quantity * s.High,
					Low:   us.Quantity * s.Low,
				}

				histMap[s.Date] = hd
			} else {

				//Pull obj from map and update values before reinserting
				hd := histMap[s.Date]
				hd.Close += (us.Quantity * s.Close)
				hd.Open += (us.Quantity * s.Open)
				hd.High += (us.Quantity * s.High)
				hd.Low += (us.Quantity * s.Low)

				histMap[s.Date] = hd
			}
		}
	}

	//Sort data and Generate History Items list. Data is sorted to help the user read output

	keys := make([]time.Time, 0, len(histMap))
	for k := range histMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	for _, key := range keys {
		hist = append(hist, histMap[key])
	}

	klogger.Exit(method)
	return hist, nil
}
