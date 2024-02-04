package polygonservice

import (
	"encoding/json"
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/internal/finance-mngr/models/restmodels"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type PolygonService struct {
	PolygonApiKey        string
	StocksEnabled        bool
	StocksApiKeyFileName string
	BaseApi string
}

// Return if stocks is enabled
func (ps *PolygonService) GetIsStocksEnabled() bool {
	method := "polygon_service.GetIsStocksEnabled"
	fmlogger.Enter(method)
	fmlogger.Exit(method)
	return ps.StocksEnabled
}

// Attempts to load an API key from a file
func (ps *PolygonService) LoadApiKeyFromFile() error {
	method := "polygon_service.LoadApiKeyFromFile"
	fmlogger.Enter(method)

	pwd, _ := os.Getwd()
	bs, err := os.ReadFile(pwd + ps.StocksApiKeyFileName)

	if err != nil {
		fmlogger.ExitError(method, "key file not found", err)
		return err
	}

	fmlogger.Info(method, "key loaded successfully")

	ps.PolygonApiKey = string(bs)
	ps.StocksEnabled = true

	fmlogger.Exit(method)
	return nil
}

// Reads an API key into the application object and persists it into a file
func (ps *PolygonService) UpdateAndPersistAPIKey(k string) error {
	method := "polygon_service.UpdateAndPersistAPIKey"
	fmlogger.Enter(method)

	fmlogger.Info(method, "Loading key into application")
	ps.PolygonApiKey = k
	ps.StocksEnabled = true

	fmlogger.Info(method, "attempting to persist API key file")
	pwd, _ := os.Getwd()

	err := os.WriteFile(pwd+ps.StocksApiKeyFileName, []byte(k), 0666)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when writing key file", err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}

func (ps *PolygonService) makeExternalCall(api string) ([]byte, error) {
	method := "polygon_service.makeExternalCall"
	fmlogger.Enter(method)

	uri := api + "?apiKey=" + ps.PolygonApiKey
	fmlogger.Info(method, "attempting to call external uri %s", api)

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

func (ps *PolygonService) FetchStockWithTicker(ticker string) (models.Stock, error) {
	method := "polygon_service.fetchStockWithTicker"
	fmlogger.Enter(method)

	api := fmt.Sprintf(ps.BaseApi+constants.PolygonGetPrevCloseAPI, ticker)
	resp, err := ps.makeExternalCall(api)
	var s models.Stock
	var pc restmodels.AggResponse
	var pci restmodels.AggResponseItem

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

func (ps *PolygonService) FetchStockWithTickerForPastYear(ticker string) ([]models.Stock, error) {
	method := "polygon_service.FetchStockWithTickerForPastYear"
	fmlogger.Enter(method)

	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	limit := time.Date(today.Year()-1, today.Month(), today.Day(), 0, 0, 0, 0, time.Local)

	api := fmt.Sprintf(ps.BaseApi+constants.PolygonGetDateRangeAPI, ticker, fmt.Sprint(limit.Format("2006-01-02")), fmt.Sprint(today.Format("2006-01-02")))
	fmt.Printf("calling external api: %s\n", api)
	resp, err := ps.makeExternalCall(api)
	var s []models.Stock
	var pc restmodels.AggResponse

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

func (ps *PolygonService) FetchStockWithTickerForDateRange(t string, d1 time.Time, d2 time.Time) ([]models.Stock, error) {
	method := "polygon_service.FetchStockWithTickerForDateRange"
	fmlogger.Enter(method)

	api := fmt.Sprintf(ps.BaseApi+constants.PolygonGetDateRangeAPI, t, fmt.Sprint(d1.Format("2006-01-02")), fmt.Sprint(d2.Format("2006-01-02")))

	fmt.Printf("[%s] attempting to call external API: %s\n", method, api)

	resp, err := ps.makeExternalCall(api)
	var s []models.Stock
	var pc restmodels.AggResponse

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
