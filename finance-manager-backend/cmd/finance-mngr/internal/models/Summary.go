package models

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"sort"
	"time"
)

const expenseType = "expense"
const incomeType = "income"
const loanSrc = "loan"
const incomeSrc = "income"
const taxSrc = "taxes"
const taxName = "income tax"
const billSrc = "bill"
const ccSrc = "credit-card"

type SummaryItem struct {
	Type    string  `json:"type"`
	Source  string  `json:"source"`
	Name    string  `json:"name"`
	Amount  float64 `json:"amount"`
	Balance float64 `json:"balance"`
}

type ExpenseSummary struct {
	Expenses          []SummaryItem `json:"expenses"`
	TotalCost         float64       `json:"totalCost"`
	TotalBalance      float64       `json:"totalBalance"`
	LoanCost          float64       `json:"loanCost"`
	LoanBalance       float64       `json:"loanBalance"`
	Taxes             float64       `json:"taxes"`
	BillCost          float64       `json:"bills"`
	CreditCardCost    float64       `json:"creditCards"`
	CreditCardBalance float64       `json:"creditCardBalance"`
	OverallBalance    float64       `json:"overallBalance"`
}

type IncomeSummary struct {
	Incomes     []SummaryItem `json:"incomes"`
	TotalIncome float64       `json:"totalIncome"`
}

type Summary struct {
	IncomeSummary  IncomeSummary  `json:"incomeSummary"`
	ExpenseSummary ExpenseSummary `json:"expenseSummary"`
	NetFunds       float64        `json:"netFunds"`
}

func (e *ExpenseSummary) CalculateExpenses() {
	method := "Summary.CalculateExpenses"
	fmlogger.Enter(method)

	e.TotalCost = e.LoanCost + e.Taxes + e.BillCost + e.CreditCardCost
	e.TotalBalance = e.LoanBalance + e.CreditCardBalance

	fmlogger.Exit(method)
}

func (s *Summary) Finalize() {
	method := "Summary.Finalize"
	fmlogger.Enter(method)

	s.NetFunds = s.IncomeSummary.TotalIncome - s.ExpenseSummary.TotalCost

	//Sort Expenses by amount
	sort.Slice(s.ExpenseSummary.Expenses, func(i, j int) bool {
		return s.ExpenseSummary.Expenses[i].Amount > s.ExpenseSummary.Expenses[j].Amount
	})

	//Sort Incomes by amount
	sort.Slice(s.IncomeSummary.Incomes, func(i, j int) bool {
		return s.IncomeSummary.Incomes[i].Amount > s.IncomeSummary.Incomes[j].Amount
	})

	fmlogger.Exit(method)
}

func (s *Summary) LoadLoans(larr []*Loan) {
	method := "Summary.loadLoans"
	fmlogger.Enter(method)

	loanBalance := 0.0
	loanCost := 0.0

	//Loop through each loan and create an item for it
	for _, l := range larr {
		i := SummaryItem{
			Type:    expenseType,
			Source:  loanSrc,
			Name:    l.Name,
			Amount:  l.MonthlyPayment,
			Balance: l.Total,
		}

		//Add new item and increment total values
		s.ExpenseSummary.Expenses = append(s.ExpenseSummary.Expenses, i)
		loanBalance += i.Balance
		loanCost += l.MonthlyPayment
	}

	//Set total values
	s.ExpenseSummary.LoanBalance = loanBalance
	s.ExpenseSummary.LoanCost = loanCost
	s.ExpenseSummary.CalculateExpenses()

	fmlogger.Enter(method)
}

func (s *Summary) LoadIncomes(iarr []*Income) {
	method := "Summary.LoadIncomes"
	fmlogger.Enter(method)

	totalIncome := 0.0
	taxes := 0.0

	//Loop through each income and add up values
	for _, i := range iarr {
		j := SummaryItem{
			Type:   incomeType,
			Source: incomeSrc,
			Name:   i.Name,
			Amount: i.GetMonthlyGrossPay(time.Now()),
		}

		//Add new item and increment total values
		s.IncomeSummary.Incomes = append(s.IncomeSummary.Incomes, j)
		totalIncome += j.Amount
		taxes += i.GetMonthlyTaxes(time.Now())
	}

	//Set Gross income for this month
	s.IncomeSummary.TotalIncome = totalIncome

	//Add Tax expense to expenses
	if taxes > 0 {
		taxItem := SummaryItem{
			Type:   expenseType,
			Source: taxSrc,
			Name:   taxName,
			Amount: taxes,
		}

		s.ExpenseSummary.Expenses = append(s.ExpenseSummary.Expenses, taxItem)
	}

	s.ExpenseSummary.Taxes = taxes
	s.ExpenseSummary.CalculateExpenses()

	fmlogger.Exit(method)
}

func (s *Summary) LoadBills(barr []*Bill) {
	method := "Summary.LoadBills"
	fmlogger.Enter(method)

	totalCost := 0.0

	//Loop through each bill and add up values
	for _, b := range barr {
		i := SummaryItem{
			Type:   expenseType,
			Source: billSrc,
			Name:   b.Name,
			Amount: b.Amount,
		}

		//Add new item and increment total values
		s.ExpenseSummary.Expenses = append(s.ExpenseSummary.Expenses, i)
		totalCost += b.Amount
	}

	//Set total cost for the month
	s.ExpenseSummary.BillCost = totalCost

	//Recalculate total cost
	s.ExpenseSummary.CalculateExpenses()

	fmlogger.Exit(method)
}

func (s *Summary) LoadCreditCards(carr []*CreditCard) {
	method := "Summary.LoadCreditCards"
	fmlogger.Enter(method)

	tc := 0.0
	tb := 0.0

	//Loop through each credit card and add up the values
	for _, cc := range carr {
		cc.CalcPayment()

		i := SummaryItem{
			Type:    expenseType,
			Source:  ccSrc,
			Name:    cc.Name,
			Amount:  cc.Payment,
			Balance: cc.Balance,
		}

		//Add new item and increment total values
		s.ExpenseSummary.Expenses = append(s.ExpenseSummary.Expenses, i)
		tc += cc.Payment
		tb += cc.Balance
	}

	//Set totals for the month
	s.ExpenseSummary.CreditCardCost = tc
	s.ExpenseSummary.CreditCardBalance = tb

	//Recalculate total cost
	s.ExpenseSummary.CalculateExpenses()

	fmlogger.Exit(method)
}
