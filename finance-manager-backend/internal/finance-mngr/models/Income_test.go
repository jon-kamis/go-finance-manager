package models

import (
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/fmUtil"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testDate = time.Date(2024, 1, 23, 7, 45, 23, 0, time.UTC)

func TestPopulateEmptyValues(t *testing.T) {
	method := "Income_test.TestPopulateEmptyValues"
	fmlogger.Enter(method)

	var i Income

	err := i.PopulateEmptyValues(testDate)
	if err == nil {
		t.Errorf("expected error to be thrown for missing Rate or Gross Pay but none was thrown")
	}

	i.Rate = 25.00
	i.TaxPercentage = -10.0

	err = i.PopulateEmptyValues(testDate)
	if err == nil {
		t.Errorf("expected error to be thrown for invalid tax percentage but none was thrown")
	}

	i.TaxPercentage = 25.0
	err = i.PopulateEmptyValues(testDate)
	if err != nil {
		t.Errorf("unexpcted error thrown for good case")
	}

	fmlogger.Exit(method)
}
func TestPopulateEmptyValues_weekly(t *testing.T) {
	method := "Income_test.TestPopulateEmptyValues_weekly"
	fmlogger.Enter(method)

	i := Income{
		Name:          "Weekly Income",
		Rate:          25.00,
		TaxPercentage: .250,
		Frequency:     constants.IncomeFreqWeekly,
		Type:          constants.IncomeTypeHourly,
		StartDt:       time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}

	//Test Calculated Values
	err := i.PopulateEmptyValues(testDate)
	if err != nil {
		t.Errorf("unexpcted error thrown for good case")
	}

	// We should see 1 pay of 40 hours
	// Expected Values:
	// NextPay: JAN 26 2024
	// GrossPay: 1000
	// Taxes: 250
	// NetPay: 750

	assert.Equal(t, time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), i.NextDt)
	assert.Equal(t, 40.0, i.Hours)
	assert.Equal(t, 1000.0, i.GrossPay)
	assert.Equal(t, 250.0, i.Taxes)
	assert.Equal(t, 750.0, i.NetPay)

	fmlogger.Exit(method)
}

func TestPopulateEmptyValues_biweekly(t *testing.T) {
	method := "Income_test.TestPopulateEmptyValues_biweekly"
	fmlogger.Enter(method)

	i := Income{
		Name:          "Bi-Weekly Income",
		Rate:          25.00,
		TaxPercentage: .250,
		Frequency:     constants.IncomeFreqBiWeekly,
		Type:          constants.IncomeTypeHourly,
		StartDt:       time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}

	//Test Calculated Values
	err := i.PopulateEmptyValues(testDate)
	if err != nil {
		t.Errorf("unexpcted error thrown for good case")
	}

	//We should see 1 pay of 80 hours
	//Expected Values:
	//NextPay: JAN 26 2024
	//GrossPay: 2000
	//Taxes: 500
	//NetPay: 1500

	assert.Equal(t, time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC), i.NextDt)
	assert.Equal(t, 80.0, i.Hours)
	assert.Equal(t, 2000.0, i.GrossPay)
	assert.Equal(t, 500.0, i.Taxes)
	assert.Equal(t, 1500.0, i.NetPay)
}

func TestPopulateEmptyValues_monthly(t *testing.T) {
	method := "Income_test.TestPopulateEmptyValues_monthly"
	fmlogger.Enter(method)

	i := Income{
		Name:          "Monthly Income",
		Rate:          6250.00,
		TaxPercentage: .250,
		Frequency:     constants.IncomeFreqMonthly,
		Type:          constants.IncomeTypeSalary,
		StartDt:       time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}

	//Test Calculated Values
	err := i.PopulateEmptyValues(testDate)
	if err != nil {
		t.Errorf("unexpcted error thrown for good case")
	}

	//We should see 1 pay of 168 hours
	//Expected Values:
	//NextPay: JAN 26 2024
	//GrossPay: 6250
	//Taxes: 1562.5
	//NetPay: 4687.5

	assert.Equal(t, time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC), i.NextDt)
	assert.Equal(t, 168.0, i.Hours)
	assert.Equal(t, 6250.0, i.GrossPay)
	assert.Equal(t, 1562.5, i.Taxes)
	assert.Equal(t, 4687.5, i.NetPay)
}

func TestGetPaysForMonthContainingDate_weekly(t *testing.T) {
	method := "Income_test.TestGetPaysForMonthContainingDate_weekly"
	fmlogger.Enter(method)

	//Testing with Feb 2024
	//There are 5 Thursdays, and 4 of every other date
	//The 1st is a Thursday

	d0 := time.Date(2024, 1, 23, 0, 0, 0, 0, time.UTC)
	d1 := time.Date(2024, 2, 23, 0, 0, 0, 0, time.UTC)

	i := Income{
		StartDt:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		Frequency: constants.IncomeFreqWeekly,
	}

	//Starts in February, This is jan so should be 0
	p := i.GetPaysForMonthContainingDate(d0)
	assert.Equal(t, 0, p)

	//5 Thursdays
	p = i.GetPaysForMonthContainingDate(d1)
	assert.Equal(t, 5, p)

	//4 Fridays
	i.StartDt = time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC)
	p = i.GetPaysForMonthContainingDate(d1)
	assert.Equal(t, 4, p)

	//Mid-Month start, only 2 paydays
	i.StartDt = time.Date(2024, 2, 21, 0, 0, 0, 0, time.UTC)
	p = i.GetPaysForMonthContainingDate(d1)
	assert.Equal(t, 2, p)

	fmlogger.Exit(method)
}

func TestGetPaysForMonthContainingDate_biweekly(t *testing.T) {
	method := "Income_test.TestGetPaysForMonthContainingDate_biweekly"
	fmlogger.Enter(method)

	d1 := time.Date(2024, 1, 23, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2024, 2, 23, 0, 0, 0, 0, time.UTC)
	d3 := time.Date(2024, 4, 23, 7, 6, 5, 4, time.UTC)

	i := Income{
		StartDt:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		Frequency: constants.IncomeFreqBiWeekly,
	}

	p := i.GetPaysForMonthContainingDate(d1)
	assert.Equal(t, 0, p)

	p = i.GetPaysForMonthContainingDate(d2)
	assert.Equal(t, 3, p)

	p = i.GetPaysForMonthContainingDate(d3)
	assert.Equal(t, 2, p)

	fmlogger.Exit(method)
}

func TestGetPaysForMonthContainingDate_monthly(t *testing.T) {
	method := "Income_test.TestGetPaysForMonthContainingDate_monthly"
	fmlogger.Enter(method)

	d1 := time.Date(2024, 1, 23, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2024, 2, 23, 0, 0, 0, 0, time.UTC)

	i := Income{
		StartDt:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		Frequency: constants.IncomeFreqMonthly,
	}

	p := i.GetPaysForMonthContainingDate(d1)
	assert.Equal(t, 0, p)

	p = i.GetPaysForMonthContainingDate(d2)
	assert.Equal(t, 1, p)

	fmlogger.Exit(method)
}

func TestGetMonthlyNetPay(t *testing.T) {
	method := "Income_test.TestGetMonthlyNetPay"
	fmlogger.Enter(method)

	d1 := time.Date(2024, 1, 23, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(2024, 2, 23, 0, 0, 0, 0, time.UTC)

	i := Income{
		StartDt:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		Frequency: constants.IncomeFreqMonthly,
		NetPay:    6250,
	}

	p := i.GetMonthlyNetPay(d1)
	assert.Equal(t, 0.0, p)

	p = i.GetMonthlyNetPay(d2)
	assert.Equal(t, 6250.0, p)

	i.Frequency = constants.IncomeFreqWeekly
	p = i.GetMonthlyNetPay(d2)
	assert.Equal(t, 31250.0, p)

	fmlogger.Exit(method)
}

func TestValidateTypeAndFrequency(t *testing.T) {
	method := "Income_test.TestValidateTypeAndFrequency"
	fmlogger.Enter(method)

	i := Income{
		Type:      constants.IncomeTypeHourly,
		Frequency: constants.IncomeFreqBiWeekly,
	}

	err := i.ValidateTypeAndFrequency()
	assert.Nil(t, err)

	i = Income{
		Type:      "",
		Frequency: "",
	}
	err = i.ValidateTypeAndFrequency()
	assert.NotNil(t, err)

	i = Income{
		Type:      "",
		Frequency: constants.IncomeFreqBiWeekly,
	}
	err = i.ValidateTypeAndFrequency()
	assert.NotNil(t, err)

	i = Income{
		Type:      constants.IncomeTypeHourly,
		Frequency: "",
	}
	err = i.ValidateTypeAndFrequency()
	assert.NotNil(t, err)

	i = Income{
		Type:      constants.IncomeTypeHourly,
		Frequency: constants.IncomeTypeHourly,
	}
	err = i.ValidateTypeAndFrequency()
	assert.NotNil(t, err)

	i = Income{
		Type:      constants.IncomeFreqBiWeekly,
		Frequency: constants.IncomeFreqBiWeekly,
	}
	err = i.ValidateTypeAndFrequency()
	assert.NotNil(t, err)

	fmlogger.Exit(method)
}

func TestValidateCanSaveIncome(t *testing.T) {
	method := "Income_test.TestValidateCanSaveIncome"
	fmlogger.Enter(method)

	var it Income
	i := Income{
		Name:          "TestValidateCanSaveIncome",
		GrossPay:      400,
		Hours:         40,
		Rate:          10.0,
		StartDt:       time.Date(2021, 2, 1, 2, 3, 4, 5, time.UTC),
		TaxPercentage: 0.15,
		UserID:        1,
		Type:          constants.IncomeTypeSalary,
		Frequency:     constants.IncomeFreqBiWeekly,
	}

	err := i.ValidateCanSaveIncome()
	assert.Nil(t, err)

	//Name is required
	it = i
	it.Name=""
	err = it.ValidateCanSaveIncome()
	assert.NotNil(t, err)

	//Gross Pay must be greater than 0
	it = i
	it.GrossPay = 0
	err = it.ValidateCanSaveIncome()
	assert.NotNil(t, err)

	//Hours must be greater than 0
	it = i
	it.Hours = 0
	err = it.ValidateCanSaveIncome()
	assert.NotNil(t, err)

	//Pay rate must be greater than 0
	it = i
	it.Rate = 0
	err = it.ValidateCanSaveIncome()
	assert.NotNil(t, err)

	//Start date must be defined
	it = i
	var d time.Time
	it.StartDt = d
	err = it.ValidateCanSaveIncome()
	assert.NotNil(t, err)

	//Tax percentage must be 0 or greater
	it = i
	it.TaxPercentage = -1.2
	err = it.ValidateCanSaveIncome()
	assert.NotNil(t, err)

	//UserId is required
	it = i
	it.UserID = 0
	err = it.ValidateCanSaveIncome()
	assert.NotNil(t, err)

	//Type and Frequency must both be valid
	it = i
	it.Frequency = "freq"
	err = it.ValidateCanSaveIncome()
	assert.NotNil(t, err)



	fmlogger.Exit(method)
}

func TestGetPayStartingDtForMonth(t *testing.T) {
	method := "Income_test.TestGetPayStartingDtForMonth"
	fmlogger.Enter(method)

	//Date to test against (only month matters)
	d := time.Date(2024, 2, 12, 20, 45, 6, 8, time.UTC)

	//Starting dates to compre
	s1 := time.Date(2023, 2, 23, 13, 6, 4, 5, time.UTC) //Later date in month but for prior year
	s2 := time.Date(2024, 1, 23, 10, 2, 4, 5, time.UTC) //Date is before month begin date
	s3 := time.Date(2024, 2, 15, 10, 2, 4, 5, time.UTC) //Date is after month begin date

	td1 := getWeeklyPayStartingDtForMonth(d, s1)
	td2 := getWeeklyPayStartingDtForMonth(d, s2)
	td3 := getWeeklyPayStartingDtForMonth(d, s3)

	//td1 should be beginning of the month of d since s1 is before that date
	//td2 should be beginning of the month of d since s2 is before that date
	//td3 should be beginning of the day of s3 since it is after the beginning of the month of d

	assert.Equal(t, fmUtil.GetMonthBeginDate(d), td1)
	assert.Equal(t, fmUtil.GetMonthBeginDate(d), td2)
	assert.Equal(t, fmUtil.GetStartOfDay(s3), td3)

	fmlogger.Exit(method)
}
