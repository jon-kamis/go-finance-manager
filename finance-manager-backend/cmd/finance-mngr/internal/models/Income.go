package models

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/constants"
	"finance-manager-backend/cmd/finance-mngr/internal/fmUtil"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"strings"
	"time"
)

// Type Income contains methods and fields related to a user's income source
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

// Function PopulateEmptyValues takes an argument of time and uses it
// to determine how much income will be generated for a user for the month containing that time.
// Values calculated are Hours, GrossPay, Taxes, and NetPay
func (i *Income) PopulateEmptyValues(t time.Time) error {
	method := "Income.PopulateEmptyValues"
	fmlogger.Enter(method)

	//Do validations first
	if i.Rate <= 0 {
		err := errors.New("rate is required")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	if i.TaxPercentage < 0 {
		err := errors.New("tax rate must be positive")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	//Determine Next Payday
	if strings.Compare(i.Frequency, constants.IncomeFreqWeekly) == 0 {
		payday := int(i.StartDt.Weekday())
		date := fmUtil.GetStartOfDay(t)

		for int(date.Weekday()) != payday {
			date = date.AddDate(0, 0, 1)
		}

		i.NextDt = date

	} else if strings.Compare(i.Frequency, constants.IncomeFreqBiWeekly) == 0 {
		date := fmUtil.GetStartOfDay(i.StartDt)
		for date.Before(t) {
			date = date.AddDate(0, 0, 14)
		}

		i.NextDt = date
	} else if strings.Compare(i.Frequency, constants.IncomeFreqMonthly) == 0 {
		date := fmUtil.GetStartOfDay(i.StartDt)
		for date.Before(t) {
			date = date.AddDate(0, 1, 0)
		}

		i.NextDt = date
	}

	// Populate Hours
	if i.Hours == 0 {
		hoursPerWeek := 40

		if strings.Compare(i.Frequency, constants.IncomeFreqWeekly) == 0 {
			i.Hours = float64(hoursPerWeek)
		} else if strings.Compare(i.Frequency, constants.IncomeFreqBiWeekly) == 0 {
			i.Hours = float64(hoursPerWeek * 2)
		} else if strings.Compare(i.Frequency, constants.IncomeFreqMonthly) == 0 {
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
	if i.Rate > 0 {
		if strings.Compare(i.Type, constants.IncomeTypeHourly) == 0 {
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

	fmlogger.Exit(method)
	return nil
}

func (i *Income) ValidateTypeAndFrequency() error {
	method := "Income.ValidateTypeAndFrequency"
	fmlogger.Enter(method)

	if i.Type == "" || i.Frequency == "" {
		err := errors.New("type and frequency are required")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	isValidType := false
	isValidFrequency := false

	for _, e := range constants.ValidTypes {
		if strings.Compare(i.Type, e) == 0 {
			isValidType = true
		}
	}

	for _, e := range constants.ValidFreq {
		if strings.Compare(i.Frequency, e) == 0 {
			isValidFrequency = true
		}
	}

	if !isValidType || !isValidFrequency {
		err := errors.New("type or frequency is invalid")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}

func (i *Income) ValidateCanSaveIncome() error {
	method := "Income.ValidateCanSaveIncome"
	fmlogger.Enter(method)

	if i.Name == "" {
		err := errors.New("cannot save income without a name")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	if i.GrossPay <= 0 {
		err := errors.New("gross pay is required")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	if i.Hours <= 0 {
		err := errors.New("hours are required")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	if i.Rate <= 0 {
		err := errors.New("pay rate is required")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	if i.StartDt.IsZero() {
		err := errors.New("start date is required")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	if i.TaxPercentage < 0 {
		err := errors.New("tax percentage is required")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	if i.UserID <= 0 {
		err := errors.New("userId is required")
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	err := i.ValidateTypeAndFrequency()
	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}

func (i *Income) GetPaysForMonthContainingDate(t time.Time) int {
	method := "Income.GetPaysForMonthContainingDate"
	fmlogger.Enter(method)

	e := fmUtil.GetMonthEndDate(t)
	if i.StartDt.After(e) {
		fmlogger.Info(method, "income begin date is after this month")
		fmlogger.Exit(method)
		return 0
	}

	if strings.Compare(i.Frequency, constants.IncomeFreqMonthly) == 0 {
		fmlogger.Exit(method)
		return 1
	} else if strings.Compare(i.Frequency, constants.IncomeFreqWeekly) == 0 {
		payday := int(i.StartDt.Weekday())
		pays := 0

		date := getWeeklyPayStartingDtForMonth(t, i.StartDt)
		monthEndDt := fmUtil.GetMonthEndDate(t)
		
		for int(date.Weekday()) != payday {
			date = date.AddDate(0, 0, 1)
		}

		for date.Before(monthEndDt) {
			pays++
			date = date.AddDate(0, 0, 7)
		}

		fmlogger.Exit(method)
		return pays
	} else {
		pays := 0
		date := i.StartDt
		monthBegin := fmUtil.GetMonthBeginDate(t)
		monthEndDt := fmUtil.GetMonthEndDate(t)

		//Iterate to the first payday of the current month
		for date.Before(monthBegin) {
			date = date.AddDate(0, 0, 14)
		}

		//Iterate and count paydays until beginning of next month
		for !date.After(monthEndDt) {
			pays++
			date = date.AddDate(0, 0, 14)
		}

		fmlogger.Exit(method)
		return pays
	}
}

func (i *Income) GetMonthlyNetPay(t time.Time) float64 {
	method := "Income.GetMonthlyNetPay"
	fmlogger.Enter(method)

	netPay := i.NetPay * float64(i.GetPaysForMonthContainingDate(t))

	fmlogger.Exit(method)
	return netPay
}

func (i *Income) GetMonthlyGrossPay(t time.Time) float64 {
	method := "Income.GetMonthlyNetPay"
	fmlogger.Enter(method)

	grossPay := i.GrossPay * float64(i.GetPaysForMonthContainingDate(t))

	fmlogger.Exit(method)
	return grossPay
}

func (i *Income) GetMonthlyTaxes(t time.Time) float64 {
	method := "Income.GetMonthlyTaxes"
	fmlogger.Enter(method)

	taxes := i.Taxes * float64(i.GetPaysForMonthContainingDate(t))

	fmlogger.Exit(method)
	return taxes
}

// Function getPayStartingDtForMonth returns either the first day of the month 't' occurs in
// or the first nanosecond of s if s is after t.
// t is the month we are getting the pay starting date for
// s is the StartDt of the income this is being called for
func getWeeklyPayStartingDtForMonth(t time.Time, s time.Time) time.Time {
	method := "Income.getPayStartingDtForMonth"
	fmlogger.Enter(method)

	monthBeginDt := fmUtil.GetMonthBeginDate(t)
	startDt := fmUtil.GetStartOfDay(s)

	if startDt.Before(monthBeginDt) {
		fmlogger.Exit(method)
		return monthBeginDt
	}

	fmlogger.Exit(method)
	return startDt
}
