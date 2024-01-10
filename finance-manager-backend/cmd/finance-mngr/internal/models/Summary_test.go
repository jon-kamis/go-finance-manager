package models

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateExpenses(t *testing.T) {
	method := "Summary_test.TestCalculateExpenses"
	fmlogger.Enter(method)

	e := ExpenseSummary{
		LoanCost:          1,
		LoanBalance:       1,
		Taxes:             1,
		BillCost:          1,
		CreditCardCost:    1,
		CreditCardBalance: 1,
	}

	e.CalculateExpenses()

	assert.Equal(t, 4.0, e.TotalCost)
	assert.Equal(t, 2.0, e.TotalBalance)

	fmlogger.Exit(method)
}

func TestFinalize(t *testing.T) {
	method := "Summary_test.TestFinalize"
	fmlogger.Enter(method)

	s := Summary{
		IncomeSummary: IncomeSummary{
			TotalIncome: 2,
		},
		ExpenseSummary: ExpenseSummary{
			TotalCost: 1,
		},
	}

	s.Finalize()
	assert.Equal(t, 1.0, s.NetFunds)

	fmlogger.Exit(method)
}

func TestLoadLoans(t *testing.T) {
	method := "Summary_test.TestLoadLoans"
	fmlogger.Enter(method)

	var s Summary
	larr := mockLoans()

	s.LoadLoans(larr)

	// All Loans initialized with Total and Payments of 1 to make test calculations easy
	assert.Equal(t, float64(len(larr)), s.ExpenseSummary.TotalBalance)
	assert.Equal(t, float64(len(larr)), s.ExpenseSummary.TotalCost)
	assert.Equal(t, float64(len(larr)), s.ExpenseSummary.LoanBalance)
	assert.Equal(t, float64(len(larr)), s.ExpenseSummary.LoanCost)

	//Expenses should have the same length as larr since only loans were added
	assert.Equal(t, len(larr), len(s.ExpenseSummary.Expenses))

	fmlogger.Exit(method)
}

func TestLoadBills(t *testing.T) {
	method := "Summary_test.TestLoadBills"
	fmlogger.Enter(method)

	var s Summary
	barr := mockBills()

	s.LoadBills(barr)

	// All Bills initialized with Total and Payments of 1 to make test calculations easy
	assert.Equal(t, 0.0, s.ExpenseSummary.TotalBalance)
	assert.Equal(t, float64(len(barr)), s.ExpenseSummary.TotalCost)
	assert.Equal(t, float64(len(barr)), s.ExpenseSummary.BillCost)

	//Expenses should have the same length as array since only it was loaded
	assert.Equal(t, len(barr), len(s.ExpenseSummary.Expenses))

	fmlogger.Exit(method)
}

func TestLoadCreditCards(t *testing.T) {
	method := "Summary_test.TestLoadCreditCards"
	fmlogger.Enter(method)

	var s Summary
	carr := mockCreditCards()

	s.LoadCreditCards(carr)

	// All Credit cards initialized with Total and Payments of 1 to make test calculations easy
	assert.Equal(t, float64(len(carr)), s.ExpenseSummary.TotalBalance)
	assert.Equal(t, float64(len(carr)), s.ExpenseSummary.CreditCardBalance)
	assert.Equal(t, float64(len(carr)), s.ExpenseSummary.TotalCost)
	assert.Equal(t, float64(len(carr)), s.ExpenseSummary.CreditCardCost)
	assert.Equal(t, float64(len(carr)), s.CreditSummary.Total)
	assert.Equal(t, 0.0, s.CreditSummary.Available)
	assert.Equal(t, 100.0, s.CreditSummary.Utilization)

	//Expenses should have the same length as array since only it was loaded
	assert.Equal(t, len(carr), len(s.ExpenseSummary.Expenses))

	fmlogger.Exit(method)
}

func TestLoadIncomes(t *testing.T) {
	method := "Summary_test.TestLoadIncomes"
	fmlogger.Enter(method)

	var s Summary
	iarr := mockIncomes()

	s.LoadIncomes(iarr)

	//All Incomes initialized where taxes and income amount should be 1 per month per income for easy calcs
	assert.Equal(t, float64(len(iarr)), s.ExpenseSummary.Taxes)
	assert.Equal(t, float64(len(iarr)), s.IncomeSummary.TotalIncome)

	//Incomes should have the same length as array since only it was loaded
	assert.Equal(t, len(iarr), len(s.IncomeSummary.Incomes))

	fmlogger.Exit(method)
}

func mockLoans() []*Loan {

	l1 := Loan{
		Name:           "Loan1",
		Total:          1,
		MonthlyPayment: 1,
	}
	l2 := Loan{
		Name:           "Loan2",
		Total:          1,
		MonthlyPayment: 1,
	}

	var larr []*Loan
	larr = append(larr, &l1)
	larr = append(larr, &l2)
	return larr

}

func mockBills() []*Bill {

	b1 := Bill{
		Name:   "Bill1",
		Amount: 1,
	}
	b2 := Bill{
		Name:   "Bill2",
		Amount: 1,
	}

	var barr []*Bill
	barr = append(barr, &b1)
	barr = append(barr, &b2)
	return barr

}

func mockCreditCards() []*CreditCard {

	c1 := CreditCard{
		Name:                 "CC1",
		Limit:                1,
		Payment:              1,
		MinPayment:           1,
		MinPaymentPercentage: 10,
		Balance:              1,
	}

	c2 := CreditCard{
		Name:                 "CC2",
		Limit:                1,
		Payment:              1,
		MinPayment:           1,
		MinPaymentPercentage: 10,
		Balance:              1,
	}

	var carr []*CreditCard
	carr = append(carr, &c1)
	carr = append(carr, &c2)
	return carr
}

func mockIncomes() []*Income {
	i1 := Income{
		Name:      "i1",
		GrossPay:  1,
		Taxes:     1,
		Frequency: "monthly",
	}

	i2 := Income{
		Name:      "i2",
		GrossPay:  1,
		Taxes:     1,
		Frequency: "monthly",
	}

	var iarr []*Income
	iarr = append(iarr, &i1)
	iarr = append(iarr, &i2)
	return iarr
}
