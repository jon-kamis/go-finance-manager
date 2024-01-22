package fmstockservice

import (
	"encoding/json"
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/internal/finance-mngr/repository"
	"finance-manager-backend/internal/finance-mngr/stockservice/fmstockservice/responsemodels"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"
)

const base_api = "https://api.polygon.io/v2"
const prev_close = "/aggs/ticker/%s/prev"
const past_year = "/aggs/ticker/%s/range/1/day/%s/%s"

type FmStockService struct {
	PolygonApiKey        string
	StocksEnabled        bool
	StocksApiKeyFileName string
	DB                   repository.DatabaseRepo
}

// Return if stocks is enabled
func (fss *FmStockService) GetIsStocksEnabled() bool {
	method := "fmstockservice.GetIsStocksEnabled"
	fmlogger.Enter(method)
	fmlogger.Exit(method)
	return fss.StocksEnabled
}

// Attempts to load an API key from a file
func (fss *FmStockService) LoadApiKeyFromFile() error {
	method := "fm_stockservice.LoadApiKeyFromFile"
	fmlogger.Enter(method)

	pwd, _ := os.Getwd()
	bs, err := os.ReadFile(pwd + fss.StocksApiKeyFileName)

	if err != nil {
		fmlogger.ExitError(method, "key file not found", err)
		return err
	}

	fmlogger.Info(method, "key loaded successfully")

	fss.PolygonApiKey = string(bs)
	fss.StocksEnabled = true

	fmlogger.Exit(method)
	return nil
}

// Reads an API key into the application object and persists it into a file
func (fss *FmStockService) UpdateAndPersistAPIKey(k string) error {
	method := "fm_stockservice.UpdateAndPersistAPIKey"
	fmlogger.Enter(method)

	fmlogger.Info(method, "Loading key into application")
	fss.PolygonApiKey = k
	fss.StocksEnabled = true

	fmlogger.Info(method, "attempting to persist API key file")
	pwd, _ := os.Getwd()

	err := os.WriteFile(pwd+fss.StocksApiKeyFileName, []byte(k), 0666)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when writing key file", err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}

func (fss *FmStockService) makeExternalCall(api string) ([]byte, error) {
	method := "fm_stockservice.makeExternalCall"
	fmlogger.Enter(method)

	uri := api + "?apiKey=" + fss.PolygonApiKey

	response, err := http.Get(uri)
	if err != nil {
		fmlogger.ExitError(method, constants.UnexpectedExternalCallError, err)
		return nil, err
	}

	if response.StatusCode != 200 {
		err = errors.New(constants.UnexpectedResponseCodeError)
		fmlogger.ExitError(method, constants.UnexpectedExternalCallError, err)
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmlogger.ExitError(method, constants.FailedToParseJsonBodyError, err)
		return nil, err
	}

	fmlogger.Exit(method)
	return responseData, nil
}

func (fss *FmStockService) FetchStockWithTicker(ticker string) (models.Stock, error) {
	method := "fm_stockservice.fetchStockWithTicker"
	fmlogger.Enter(method)

	api := fmt.Sprintf(base_api+prev_close, ticker)
	resp, err := fss.makeExternalCall(api)
	var s models.Stock
	var pc responsemodels.AggResponse
	var pci responsemodels.AggResponseItem

	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		return s, err
	}

	err = json.Unmarshal(resp, &pc)
	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		return s, err
	}

	pci = pc.Results[0]

	s = models.Stock{
		Ticker: pc.Ticker,
		High:   pci.High,
		Low:    pci.Low,
		Open:   pci.Open,
		Close:  pci.Close,
		Date:   time.UnixMilli(int64(pci.UnixTime)),
	}

	fmlogger.Exit(method)
	return s, nil
}

func (fss *FmStockService) FetchStockWithTickerForPastYear(ticker string) ([]models.Stock, error) {
	method := "fm_stockservice.FetchStockWithTickerForPastYear"
	fmlogger.Enter(method)

	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	limit := time.Date(today.Year()-1, today.Month(), today.Day(), 0, 0, 0, 0, time.Local)

	api := fmt.Sprintf(base_api+past_year, ticker, fmt.Sprint(limit.Format("2006-01-02")), fmt.Sprint(today.Format("2006-01-02")))
	fmt.Printf("calling external api: %s\n", api)
	resp, err := fss.makeExternalCall(api)
	var s []models.Stock
	var pc responsemodels.AggResponse

	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		return s, err
	}

	err = json.Unmarshal(resp, &pc)
	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		return s, err
	}

	for _, pci := range pc.Results {
		i := models.Stock{
			Ticker: pc.Ticker,
			High:   pci.High,
			Low:    pci.Low,
			Open:   pci.Open,
			Close:  pci.Close,
			Date:   time.UnixMilli(int64(pci.UnixTime)),
		}
		s = append(s, i)
	}

	fmlogger.Exit(method)
	return s, nil
}

func (fss *FmStockService) FetchStockWithTickerForDateRange(t string, d1 time.Time, d2 time.Time) ([]models.Stock, error) {
	method := "fm_stockservice.FetchStockWithTickerForDateRange"
	fmlogger.Enter(method)

	api := fmt.Sprintf(base_api+past_year, t, fmt.Sprint(d1.Format("2006-01-02")), fmt.Sprint(d2.Format("2006-01-02")))

	resp, err := fss.makeExternalCall(api)
	var s []models.Stock
	var pc responsemodels.AggResponse

	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		return s, err
	}

	err = json.Unmarshal(resp, &pc)
	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		return s, err
	}

	for _, pci := range pc.Results {
		i := models.Stock{
			Ticker: pc.Ticker,
			High:   pci.High,
			Low:    pci.Low,
			Open:   pci.Open,
			Close:  pci.Close,
			Date:   time.UnixMilli(int64(pci.UnixTime)),
		}
		s = append(s, i)
	}

	fmlogger.Exit(method)
	return s, nil
}

// Function GetUserPortfolioBalanceHistory fetches the Portfolio Balance History for a user for a given timeframe
// uId - The ID of the user to fetch history for
// d - The number of past days to pull history for. Maximum is 365
func (fss *FmStockService) GetUserPortfolioBalanceHistory(uId int, d int) ([]models.PortfolioBalanceHistory, error) {
	method := "fm_stockservice.GetUserPortfolioBalanceHistory"
	fmlogger.Enter(method)

	var hist []models.PortfolioBalanceHistory
	var sd time.Time
	var err error
	ed := time.Now()

	//Validate passed values
	if uId == 0 {
		err = errors.New("uId is required")
		fmlogger.ExitError(method, err.Error(), err)
		return hist, err
	}

	if d < 1 || d > 365 {
		err = errors.New("d must be between 1 and 365 inclusively")
		fmlogger.ExitError(method, err.Error(), err)
		return hist, err
	}

	//First Load User Positions for date range
	sd = ed.Add(-1 * time.Duration(d) * 24 * time.Hour)

	usl, err := fss.DB.GetAllUserStocksByDateRange(uId, "", sd, ed)

	if err != nil {
		fmlogger.ExitError(method, constants.UnexpectedSQLError, err)
		return hist, err
	}

	//No stocks exist
	if len(usl) == 0 {
		fmlogger.Info(method, "User has no stocks for this timeframe")
		fmlogger.Exit(method)
		return hist, nil
	}

	//Next Loop through each user position and load stock data for that position. Add total value for each date
	histMap := make(map[time.Time]float64)

	for _, us := range usl {

		var d1 time.Time
		var d2 time.Time

		if us.EffectiveDt.Before(sd) {
			d1 = sd
		} else {
			d1 = us.EffectiveDt
		}

		if !us.ExpirationDt.IsZero() && us.ExpirationDt.Before(ed) {
			d2 = us.ExpirationDt
		} else {
			d2 = ed
		}

		sl, err := fss.DB.GetStockDataByTickerAndDateRange(us.Ticker, d1, d2)

		if err != nil {
			fmlogger.ExitError(method, constants.UnexpectedSQLError, err)
			return hist, err
		}

		//Next, Loop through stock Data for this entry and add totals to each date in map
		for _, s := range sl {
			val := (us.Quantity * s.Close)
			if histMap[s.Date] == 0 {
				histMap[s.Date] = val
			} else {
				histMap[s.Date] = val + histMap[s.Date]
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
		hist = append(hist, models.PortfolioBalanceHistory{Date: key, Balance: histMap[key]})
	}

	fmlogger.Exit(method)
	return hist, nil
}
