package restmodels

import (
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/enums/stockoperation"
	"time"

	"github.com/jon-kamis/klogger"
)

// Type ModifyStockRequest holds data for a PATCH modify stocks request
type ModifyStockRequest struct {
	//The ticker to modify
	Ticker string `json:"ticker"`

	//The amount to modify
	Amount float64 `json:"amount"`

	//The operation. Options are 'buy' and 'sell'
	Operation stockoperation.ModifyStockOperation `json:"operation"`

	//Date of operation
	Date time.Time `json:"date"`
}

func (m *ModifyStockRequest) IsValidRequest() (bool, string) {
	method := "ModifyStockRequest.isValidRequest"
	klogger.Enter(method)

	isValid := true
	var msg string

	if m.Ticker == "" {
		msg = constants.StockOperationTickerRequiredError
		isValid = false
	}

	if m.Amount <= 0 {
		msg = constants.StockOperationInvalidAmountError
		isValid = false
	}

	if m.Operation != stockoperation.Add && m.Operation != stockoperation.Remove {
		msg = constants.StockOperationInvalidOperationError
		isValid = false
	}

	if m.Date.IsZero() || m.Date.After(time.Now()) {
		msg = constants.StockOperationInvalidDateError
		isValid = false
	}

	klogger.Exit(method)
	return isValid, msg
}
