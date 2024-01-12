package fmstockservice

import (
	"encoding/json"
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/constants"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"finance-manager-backend/cmd/finance-mngr/internal/stockservice.go/fmstockservice/responsemodels"
	"fmt"
	"io"
	"net/http"
	"time"
)

const base_api = "https://api.polygon.io/v2"
const prev_close = "/aggs/ticker/%s/prev"

type FmStockService struct {
}

func makeExternalCall(api string, token string) ([]byte, error) {
	method := "fm_stockservice.makeExternalCall"
	fmlogger.Enter(method)

	uri := api + "?apiKey=" + token

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

func (fss *FmStockService) FetchStockWithTicker(ticker string, token string) (models.Stock, error) {
	method := "fm_stockservice.fetchStockWithTicker"
	fmlogger.Enter(method)

	api := fmt.Sprintf(base_api+prev_close, ticker)
	resp, err := makeExternalCall(api, token)
	var s models.Stock
	var pc responsemodels.PreviousClose
	var pci responsemodels.PreviousCloseItem

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
