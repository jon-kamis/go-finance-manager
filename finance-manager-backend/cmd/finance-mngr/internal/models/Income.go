package models

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/fmUtil"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"fmt"
	"strings"
	"time"
)

var validTypes = []string{"hourly", "salary"}
var validFreq = []string{"weekly", "bi-weekly", "monthly"}

type Income struct {
	ID            int       `json:"id"`
	UserID        int       `json:"userId"`
	Name          string    `json:"name"`
	Rate          float64   `json:"rate"`
	Hours         float64   `json:"hours"`
	Type          string    `json:"type"`
	GrossPay      float64   `json:"grossPay"`
	Taxes         float64   `json:"taxes"`
	NetPay        float64   `json:"netPay"`
	Frequency     string    `json:"frequency"`
	TaxPercentage float64   `json:"taxPercentage"`
	StartDt       time.Time `json:"startDt"`
	NextDt        time.Time `json:"nextDt"`
	CreateDt      time.Time `json:"createDt"`
	LastUpdateDt  time.Time `json:"lastUpdateDt"`
}

func (i *Income) PopulateEmptyValues() error {
	method := "Income.PopulateEmptyValues"
	fmt.Printf("[ENTER %s]\n", method)

	//Do validations first
	if i.Rate <= 0 && i.GrossPay <= 0 {
		returnError("Rate or GrossPay is required", method)
	}

	if i.TaxPercentage < 0 {
		returnError("Tax rate must be positive", method)
	}

	//Determine Next Payday
	if strings.Compare(i.Frequency, "weekly") == 0 {
		payday := int(i.StartDt.Weekday())
		date := time.Now()

		for int(date.Weekday()) != payday {
			date = date.AddDate(0, 0, 1)
		}

		i.NextDt = date

	} else if strings.Compare(i.Frequency, "bi-weekly") == 0 {
		date := i.StartDt
		for date.Before(time.Now()) {
			date = date.AddDate(0, 0, 14)
		}

		i.NextDt = date
	} else if strings.Compare(i.Frequency, "monthly") == 0 {
		date := i.StartDt
		for date.Before(time.Now()) {
			date = date.AddDate(0, 1, 0)
		}

		i.NextDt = date
	}

	// Populate Hours
	if i.Hours == 0 {
		hoursPerWeek := 40

		if strings.Compare(i.Frequency, "weekly") == 0 {
			i.Hours = float64(hoursPerWeek)
		} else if strings.Compare(i.Frequency, "bi-weekly") == 0 {
			i.Hours = float64(hoursPerWeek * 2)
		} else if strings.Compare(i.Frequency, "monthly") == 0 {
			workdays := 0
			date := i.NextDt
			date = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
			nextMonth := date.AddDate(0, 1, 0)

			for date.Before(nextMonth) {
				if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
					workdays++
				}

				date = date.AddDate(0, 0, 1)
			}
			i.Hours = float64(workdays * 8)
		}
	}

	// Populate GrossPay
	if i.Rate > 0 && i.GrossPay == 0 {
		if strings.Compare(i.Type, "hourly") == 0 {
			i.GrossPay = i.Rate * i.Hours
		} else {
			i.GrossPay = i.Rate
		}
	}

	// Populate Taxes and Net Pay
	if i.Taxes == 0 {
		i.Taxes = i.GrossPay * i.TaxPercentage
		i.NetPay = i.GrossPay - i.Taxes
	}

	fmt.Printf("[EXIT %s]\n", method)
	return nil
}

func (i *Income) ValidateTypeAndFrequency() error {
	method := "Income.ValidateTypeAndFrequency"
	fmt.Printf("[ENTER %s]\n", method)

	if i.Type == "" || i.Frequency == "" {
		fmt.Printf("[%s] type and frequency are required\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		return errors.New("unprocessable request")
	}

	isValidType := false
	isValidFrequency := false

	for _, e := range validTypes {
		if strings.Compare(i.Type, e) == 0 {
			isValidType = true
		}
	}

	for _, e := range validFreq {
		if strings.Compare(i.Frequency, e) == 0 {
			isValidFrequency = true
		}
	}

	if !isValidType || !isValidFrequency {
		fmt.Printf("[%s] type or frequency is invalid\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		return errors.New("unprocessable request")
	}

	fmt.Printf("[EXIT %s]\n", method)
	return nil
}

func (i *Income) ValidateCanSaveIncome() error {
	method := "Income.ValidateCanSaveIncome"
	fmt.Printf("[ENTER %s]\n", method)

	if i.Name == "" {
		returnError("cannot save loan without a name", method)
	}

	if i.GrossPay <= 0 {
		returnError("Gross Pay is required", method)
	}

	if i.Hours <= 0 {
		returnError("Hours are required", method)
	}

	if i.Rate <= 0 {
		returnError("Pay rate is required", method)
	}

	if i.StartDt.IsZero() {
		returnError("Start date is required", method)
	}

	if i.TaxPercentage <= 0 {
		returnError("Tax rate is required", method)
	}

	if i.UserID <= 0 {
		returnError("UserId is required", method)
	}

	err := i.ValidateTypeAndFrequency()
	if err != nil {
		return err
	}

	fmt.Printf("[EXIT %s]\n", method)
	return nil
}

func (i *Income) GetPaysThisMonth() int {
	method := "Income.GetPaysThisMonth"
	fmlogger.Enter(method)

	if strings.Compare(i.Frequency, "monthly") == 0 {
		fmlogger.Exit(method)
		return 1
	} else if strings.Compare(i.Frequency, "weekly") == 0 {
		payday := int(i.StartDt.Weekday())
		pays := 0

		monthBegin := fmUtil.GetMonthBeginDate(time.Now())
		date := monthBegin
		nextMonthBegin := monthBegin.AddDate(0, 1, 0)

		date = monthBegin
		for int(date.Weekday()) != payday {
			date = date.AddDate(0, 0, 1)
		}

		for date.Before(nextMonthBegin) {
			pays++
			date = date.AddDate(0, 0, 7)
		}

		fmlogger.Exit(method)
		return pays
	} else {
		pays := 0
		date := i.StartDt
		monthBegin := fmUtil.GetMonthBeginDate(time.Now())
		nextMonthBegin := monthBegin.AddDate(0, 1, 0)

		//Iterate to the first payday of the current month
		for date.Before(monthBegin) {
			date = date.AddDate(0, 0, 14)
		}

		//Iterate and count paydays until beginning of next month
		for date.Before(nextMonthBegin) {
			pays++
			date = date.AddDate(0, 0, 14)
		}

		fmlogger.Exit(method)
		return pays
	}
}

func (i *Income) GetMonthlyNetPay() float64 {
	method := "Income.GetMonthlyNetPay"
	fmlogger.Enter(method)
	fmlogger.Exit(method)
	return i.NetPay * float64(i.GetPaysThisMonth())
}

func (i *Income) GetMonthlyGrossPay() float64 {
	method := "Income.GetMonthlyNetPay"
	fmlogger.Enter(method)
	fmlogger.Exit(method)
	return i.GrossPay * float64(i.GetPaysThisMonth())
}

func (i *Income) GetMonthlyTaxes() float64 {
	method := "Income.GetMonthlyTaxes"
	fmlogger.Enter(method)
	fmlogger.Exit(method)
	return i.Taxes * float64(i.GetPaysThisMonth())
}

func returnError(msg string, method string) error {
	fmt.Printf("[%s] %s\n", method, msg)
	fmt.Printf("[EXIT %s]\n", method)
	return errors.New(msg)
}
