package fmstockservice

import (
	"encoding/json"
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/internal/finance-mngr/stockservice.go/fmstockservice/responsemodels"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const base_api = "https://api.polygon.io/v2"
const prev_close = "/aggs/ticker/%s/prev"
const past_year = "/aggs/ticker/%s/range/1/day/%s/%s"

type FmStockService struct {
	PolygonApiKey        string
	StocksEnabled        bool
	StocksApiKeyFileName string
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
