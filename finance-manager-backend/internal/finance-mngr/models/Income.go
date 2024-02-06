package models

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/fmUtil"
	"strings"
	"time"

	"github.com/jon-kamis/klogger"
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
	klogger.Enter(method)

	//Do validations first
	if i.Rate <= 0 {
		err := errors.New("rate is required")
		klogger.ExitError(method, err.Error())
		return err
	}

	if i.TaxPercentage < 0 {
		err := errors.New("tax rate must be positive")
		klogger.ExitError(method, err.Error())
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

	klogger.Exit(method)
	return nil
}

func (i *Income) ValidateTypeAndFrequency() error {
	method := "Income.ValidateTypeAndFrequency"
	klogger.Enter(method)

	if i.Type == "" || i.Frequency == "" {
		err := errors.New("type and frequency are required")
		klogger.ExitError(method, err.Error())
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
		klogger.ExitError(method, err.Error())
		return err
	}

	klogger.Exit(method)
	return nil
}

func (i *Income) ValidateCanSaveIncome() error {
	method := "Income.ValidateCanSaveIncome"
	klogger.Enter(method)

	if i.Name == "" {
		err := errors.New("cannot save income without a name")
		klogger.ExitError(method, err.Error())
		return err
	}

	if i.GrossPay <= 0 {
		err := errors.New("gross pay is required")
		klogger.ExitError(method, err.Error())
		return err
	}

	if i.Hours <= 0 {
		err := errors.New("hours are required")
		klogger.ExitError(method, err.Error())
		return err
	}

	if i.Rate <= 0 {
		err := errors.New("pay rate is required")
		klogger.ExitError(method, err.Error())
		return err
	}

	if i.StartDt.IsZero() {
		err := errors.New("start date is required")
		klogger.ExitError(method, err.Error())
		return err
	}

	if i.TaxPercentage < 0 {
		err := errors.New("tax percentage is required")
		klogger.ExitError(method, err.Error())
		return err
	}

	if i.UserID <= 0 {
		err := errors.New("userId is required")
		klogger.ExitError(method, err.Error())
		return err
	}

	err := i.ValidateTypeAndFrequency()
	if err != nil {
		klogger.ExitError(method, err.Error())
		return err
	}

	klogger.Exit(method)
	return nil
}

func (i *Income) GetPaysForMonthContainingDate(t time.Time) int {
	method := "Income.GetPaysForMonthContainingDate"
	klogger.Enter(method)

	e := fmUtil.GetMonthEndDate(t)
	if i.StartDt.After(e) {
		klogger.Info(method, "income begin date is after this month")
		klogger.Exit(method)
		return 0
	}

	if strings.Compare(i.Frequency, constants.IncomeFreqMonthly) == 0 {
		klogger.Exit(method)
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

		klogger.Exit(method)
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

		klogger.Exit(method)
		return pays
	}
}

func (i *Income) GetMonthlyNetPay(t time.Time) float64 {
	method := "Income.GetMonthlyNetPay"
	klogger.Enter(method)

	netPay := i.NetPay * float64(i.GetPaysForMonthContainingDate(t))

	klogger.Exit(method)
	return netPay
}

func (i *Income) GetMonthlyGrossPay(t time.Time) float64 {
	method := "Income.GetMonthlyNetPay"
	klogger.Enter(method)

	grossPay := i.GrossPay * float64(i.GetPaysForMonthContainingDate(t))

	klogger.Exit(method)
	return grossPay
}

func (i *Income) GetMonthlyTaxes(t time.Time) float64 {
	method := "Income.GetMonthlyTaxes"
	klogger.Enter(method)

	taxes := i.Taxes * float64(i.GetPaysForMonthContainingDate(t))

	klogger.Exit(method)
	return taxes
}

// Function getPayStartingDtForMonth returns either the first day of the month 't' occurs in
// or the first nanosecond of s if s is after t.
// t is the month we are getting the pay starting date for
// s is the StartDt of the income this is being called for
func getWeeklyPayStartingDtForMonth(t time.Time, s time.Time) time.Time {
	method := "Income.getPayStartingDtForMonth"
	klogger.Enter(method)

	monthBeginDt := fmUtil.GetMonthBeginDate(t)
	startDt := fmUtil.GetStartOfDay(s)

	if startDt.Before(monthBeginDt) {
		klogger.Exit(method)
		return monthBeginDt
	}

	klogger.Exit(method)
	return startDt
}
