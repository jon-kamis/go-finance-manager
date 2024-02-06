package fmservice

import (
	"database/sql"
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/enums/stockoperation"
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/internal/finance-mngr/models/restmodels"
	"time"

	"github.com/jon-kamis/klogger"
)

func (fms *FMService) LoadPriorUserStockForTransaction(r restmodels.ModifyStockRequest, usp *models.UserStock, us *models.UserStock) error {
	method := "stock_service.LoadPriorUserStockForTransaction"
	klogger.Enter(method)

	var err error

	//Check for existing user stock of this ticker
	*usp, err = fms.DB.GetUserStockByUserIdTickerAndDate(us.UserId, r.Ticker, r.Date)

	if err != nil {
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return err
	}

	if usp.ID > 0 {
		//Update UserStock Before
		us.ExpirationDt = usp.ExpirationDt
		usp.ExpirationDt = sql.NullTime{Time: r.Date.Add(-1 * time.Millisecond), Valid: true}

		//Set the quantity for the new UserStock object
		switch r.Operation {
		case stockoperation.Add:
			us.Quantity = usp.Quantity + r.Amount
		case stockoperation.Remove:
			us.Quantity = usp.Quantity - r.Amount

			if us.Quantity < 0 {
				//Quantity cannot be reduced below 0
				err := errors.New(constants.StockOperationBelowZeroError)
				klogger.ExitError(method, err.Error())
				return err
			}
		}

	} else {

		if r.Operation == stockoperation.Remove {
			err := errors.New(constants.StockOperationBelowZeroError)
			klogger.ExitError(method, err.Error())
			return err
		}

		us.Quantity = r.Amount

	}

	klogger.Exit(method)
	return nil
}
