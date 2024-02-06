package polygonservice

import (
	"encoding/json"
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/internal/finance-mngr/models/restmodels"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/jon-kamis/klogger"
)

type PolygonService struct {
	PolygonApiKey        string
	StocksEnabled        bool
	StocksApiKeyFileName string
	BaseApi              string
}

// Return if stocks is enabled
func (ps *PolygonService) GetIsStocksEnabled() bool {
	method := "polygon_service.GetIsStocksEnabled"
	klogger.Enter(method)
	klogger.Exit(method)
	return ps.StocksEnabled
}

// Attempts to load an API key from a file
func (ps *PolygonService) LoadApiKeyFromFile() error {
	method := "polygon_service.LoadApiKeyFromFile"
	klogger.Enter(method)

	pwd, _ := os.Getwd()
	bs, err := os.ReadFile(pwd + ps.StocksApiKeyFileName)

	if err != nil {
		klogger.ExitError(method, "key file not found:\n%v", err)
		return err
	}

	klogger.Info(method, "key loaded successfully")

	ps.PolygonApiKey = string(bs)
	ps.StocksEnabled = true

	klogger.Exit(method)
	return nil
}

// Reads an API key into the application object and persists it into a file
func (ps *PolygonService) UpdateAndPersistAPIKey(k string) error {
	method := "polygon_service.UpdateAndPersistAPIKey"
	klogger.Enter(method)

	klogger.Info(method, "Loading key into application")
	ps.PolygonApiKey = k
	ps.StocksEnabled = true

	klogger.Info(method, "attempting to persist API key file")
	pwd, _ := os.Getwd()

	err := os.WriteFile(pwd+ps.StocksApiKeyFileName, []byte(k), 0666)

	if err != nil {
		klogger.ExitError(method, "unexpected error occured when writing key file:\n%v", err)
		return err
	}

	klogger.Exit(method)
	return nil
}

func (ps *PolygonService) makeExternalCall(api string) ([]byte, error) {
	method := "polygon_service.makeExternalCall"
	klogger.Enter(method)

	uri := api + "?apiKey=" + ps.PolygonApiKey
	klogger.Info(method, "attempting to call external uri %s", api)

	response, err := http.Get(uri)
	if err != nil {
		klogger.ExitError(method, constants.UnexpectedExternalCallError, err)
		return nil, err
	}

	if response.StatusCode != 200 {
		err = errors.New(constants.UnexpectedResponseCodeError)
		klogger.ExitError(method, constants.UnexpectedExternalCallError, err)
		return nil, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		klogger.ExitError(method, constants.FailedToParseJsonBodyError, err)
		return nil, err
	}

	klogger.Exit(method)
	return responseData, nil
}

func (ps *PolygonService) FetchStockWithTicker(ticker string) (models.Stock, error) {
	method := "polygon_service.fetchStockWithTicker"
	klogger.Enter(method)

	api := fmt.Sprintf(ps.BaseApi+constants.PolygonGetPrevCloseAPI, ticker)
	resp, err := ps.makeExternalCall(api)
	var s models.Stock
	var pc restmodels.AggResponse
	var pci restmodels.AggResponseItem

	if err != nil {
		klogger.ExitError(method, err.Error())
		return s, err
	}

	err = json.Unmarshal(resp, &pc)
	if err != nil {
		klogger.ExitError(method, err.Error())
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

	klogger.Exit(method)
	return s, nil
}

func (ps *PolygonService) FetchStockWithTickerForPastYear(ticker string) ([]models.Stock, error) {
	method := "polygon_service.FetchStockWithTickerForPastYear"
	klogger.Enter(method)

	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	limit := time.Date(today.Year()-1, today.Month(), today.Day(), 0, 0, 0, 0, time.Local)

	api := fmt.Sprintf(ps.BaseApi+constants.PolygonGetDateRangeAPI, ticker, fmt.Sprint(limit.Format("2006-01-02")), fmt.Sprint(today.Format("2006-01-02")))
	fmt.Printf("calling external api: %s\n", api)
	resp, err := ps.makeExternalCall(api)
	var s []models.Stock
	var pc restmodels.AggResponse

	if err != nil {
		klogger.ExitError(method, err.Error())
		return s, err
	}

	err = json.Unmarshal(resp, &pc)
	if err != nil {
		klogger.ExitError(method, err.Error())
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

	klogger.Exit(method)
	return s, nil
}

func (ps *PolygonService) FetchStockWithTickerForDateRange(t string, d1 time.Time, d2 time.Time) ([]models.Stock, error) {
	method := "polygon_service.FetchStockWithTickerForDateRange"
	klogger.Enter(method)

	api := fmt.Sprintf(ps.BaseApi+constants.PolygonGetDateRangeAPI, t, fmt.Sprint(d1.Format("2006-01-02")), fmt.Sprint(d2.Format("2006-01-02")))

	fmt.Printf("[%s] attempting to call external API: %s\n", method, api)

	resp, err := ps.makeExternalCall(api)
	var s []models.Stock
	var pc restmodels.AggResponse

	if err != nil {
		klogger.ExitError(method, err.Error())
		return s, err
	}

	err = json.Unmarshal(resp, &pc)
	if err != nil {
		klogger.ExitError(method, err.Error())
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

	klogger.Exit(method)
	return s, nil
}
