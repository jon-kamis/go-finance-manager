package models

import (
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPerformPaymentCalc(t *testing.T) {
	method := "Loan_test.TestPerformPaymentCalc"
	fmlogger.Enter(method)

	expectedPayment := 184.17
	l := Loan{
		Total:        10000,
		InterestRate: 4,
		LoanTerm:     60,
	}

	err := l.PerformPaymentCalc()
	assert.Nil(t, err)
	assert.Equal(t, expectedPayment, math.Round(l.MonthlyPayment*100)/100)

	l.Total = 0
	err = l.PerformPaymentCalc()
	assert.NotNil(t, err)

	fmlogger.Exit(method)
}

func TestPerformCalc(t *testing.T) {
	method := "Loan_test.TestPerformCalc"
	fmlogger.Enter(method)

	expectedInterest := 1049.30
	expectedTotalCost := 11049.30
	expectedTotalPayment := 10000.0

	l := Loan{
		Total:        0,
		InterestRate: 4,
		LoanTerm:     60,
	}

	//Invalid Request because Total is 0
	err := l.PerformCalc()
	assert.NotNil(t, err)

	l.Total = 10000
	err = l.PerformCalc()
	assert.Nil(t, err)
	assert.Equal(t, expectedInterest, math.Round(l.Interest*100)/100)
	assert.Equal(t, expectedTotalCost, math.Round(l.TotalCost*100)/100)
	assert.Equal(t, expectedTotalPayment, math.Round(l.TotalPayment*100)/100)
	assert.Equal(t, l.LoanTerm, len(l.PaymentSchedule))

	fmlogger.Exit(method)
}

func TestValidateCanPerformCalc(t *testing.T) {
	method := "Loan_test.TestPerformCalc"
	fmlogger.Enter(method)

	var lt Loan
	l := Loan{
		Total:        10000,
		InterestRate: 4,
		LoanTerm:     60,
	}

	err := l.ValidateCanPerformCalc()
	assert.Nil(t, err)

	//Loan balance must be greater than 0
	lt = l
	lt.Total = 0
	err = lt.ValidateCanPerformCalc()
	assert.NotNil(t, err)

	//Interest Rate must be 0 or greater
	lt = l
	lt.InterestRate = -1
	err = lt.ValidateCanPerformCalc()
	assert.NotNil(t, err)

	//Loan Term must be at least one month
	lt = l
	lt.LoanTerm = 0
	err = lt.ValidateCanPerformCalc()
	assert.NotNil(t, err)

	fmlogger.Exit(method)
}

func TestValidateCanSaveLoan(t *testing.T) {
	method := "Loan_test.TestPerformCalc"
	fmlogger.Enter(method)

	var lt Loan
	l := Loan{
		Name:         "TestValidateCanSaveLoan",
		Total:        10000,
		InterestRate: 4,
		LoanTerm:     60,
		UserID:       1,
		TotalCost:    11049.30,
		TotalPayment: 10000,
		Interest:     1049.91,
	}

	//Good case
	err := l.ValidateCanSaveLoan()
	assert.Nil(t, err)

	//Name is required
	lt = l
	lt.Name = ""
	err = lt.ValidateCanSaveLoan()
	assert.NotNil(t, err)

	//Total is required
	lt = l
	lt.Total = 0
	err = lt.ValidateCanSaveLoan()
	assert.NotNil(t, err)

	//Interest Rate must be >= 0
	lt = l
	lt.InterestRate = -1
	err = lt.ValidateCanSaveLoan()
	assert.NotNil(t, err)

	//Loan term is required
	lt = l
	lt.LoanTerm = 0
	err = lt.ValidateCanSaveLoan()
	assert.NotNil(t, err)

	//UserId is required
	lt = l
	lt.UserID = 0
	err = lt.ValidateCanSaveLoan()
	assert.NotNil(t, err)

	//Total Cost is required
	lt = l
	lt.TotalCost = 0
	err = lt.ValidateCanSaveLoan()
	assert.NotNil(t, err)

	//Total Payment is required
	lt = l
	lt.TotalPayment = 0
	err = lt.ValidateCanSaveLoan()
	assert.NotNil(t, err)

	//Interest must be >= 0
	lt = l
	lt.Interest = -1
	err = lt.ValidateCanSaveLoan()
	assert.NotNil(t, err)

	fmlogger.Exit(method)
}
