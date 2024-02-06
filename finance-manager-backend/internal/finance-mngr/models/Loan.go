package models

import (
	"errors"
	"math"
	"time"

	"github.com/jon-kamis/klogger"
)

type Loan struct {
	ID              int                   `json:"id"`
	UserID          int                   `json:"userId"`
	Name            string                `json:"name"`
	Total           float64               `json:"total"`
	InterestRate    float64               `json:"interestRate"`
	MonthlyPayment  float64               `json:"monthlyPayment"`
	Interest        float64               `json:"interest"`
	TotalCost       float64               `json:"totalCost"`
	TotalPayment    float64               `json:"totalPayment"`
	LoanTerm        int                   `json:"loanTerm"`
	PaymentSchedule []PaymentScheduleItem `json:"paymentSchedule"`
	CreateDt        time.Time             `json:"-"`
	LastUpdateDt    time.Time             `json:"-"`
}

type PaymentScheduleItem struct {
	Month            int     `json:"month"`
	Principal        float64 `json:"principal"`
	Interest         float64 `json:"interest"`
	InterestToDate   float64 `json:"interestToDate"`
	PrincipalToDate  float64 `json:"principalToDate"`
	RemainingBalance float64 `json:"remainingBalance"`
}

type PaymentScheduleComparisonItem struct {
	Month                 int     `json:"month"`
	Principal             float64 `json:"principal"`
	PrincipalNew          float64 `json:"principalNew"`
	PrincipalDelta        float64 `json:"principalDelta"`
	Interest              float64 `json:"interest"`
	InterestNew           float64 `json:"interestNew"`
	InterestDelta         float64 `json:"interestDelta"`
	InterestToDate        float64 `json:"interestToDate"`
	InterestToDateNew     float64 `json:"interestToDateNew"`
	InterestToDateDelta   float64 `json:"interestToDateDelta"`
	PrincipalToDate       float64 `json:"principalToDate"`
	PrincipalToDateNew    float64 `json:"principalToDateNew"`
	PrincipalToDateDelta  float64 `json:"principalToDateDelta"`
	RemainingBalance      float64 `json:"remainingBalance"`
	RemainingBalanceNew   float64 `json:"remainingBalanceNew"`
	RemainingBalanceDelta float64 `json:"remainingBalanceDelta"`
}

type LoansSummary struct {
	Count        int     `json:"count"`
	TotalBalance float64 `json:"totalBalance"`
	MonthlyCost  float64 `json:"monthlyCost"`
}

func (l *Loan) PerformCalc() error {
	method := "Loan.performCalc"
	klogger.Enter(method)

	totalCalc := l.Total
	months := 0
	var totalInterest float64
	var totalPayment float64
	var interestToDate float64
	var principalToDate float64
	var paymentSchedule []PaymentScheduleItem
	totalInterest = 0
	interestToDate = 0
	principalToDate = 0
	totalPayment = 0

	// Calculate monthly payment
	if l.MonthlyPayment < 1 {
		err := l.PerformPaymentCalc()

		if err != nil {
			errMsg := "Loan calculation request is invalid:\n%s"
			klogger.ExitError(method, errMsg, err)
			return err
		}
	}

	for totalCalc > 0 {
		interest := (totalCalc * (l.InterestRate / 100)) / 12
		thisPay := l.MonthlyPayment - interest

		interestToDate += interest
		principalToDate += thisPay
		months++
		if totalCalc-thisPay > 0.009 {
			totalCalc -= thisPay
		} else {
			thisPay = totalCalc - interest
			totalCalc = 0.0
		}
		totalInterest = totalInterest + interest
		totalPayment = totalPayment + thisPay

		paymentSum := PaymentScheduleItem{
			Month:            months,
			Principal:        thisPay,
			Interest:         interest,
			InterestToDate:   interestToDate,
			PrincipalToDate:  principalToDate,
			RemainingBalance: totalCalc,
		}

		paymentSchedule = append(paymentSchedule, paymentSum)
	}

	calcErr := l.Total - totalPayment
	klogger.Info(method, "Total: %f, totalPayment: %f, totalInterest: %f, Calculation Err: %f\n", l.Total, totalPayment, totalInterest, calcErr)

	totalPayment += calcErr
	totalInterest -= calcErr

	l.TotalPayment = totalPayment
	l.Interest = totalInterest
	l.TotalCost = totalPayment + totalInterest
	l.LoanTerm = months
	l.PaymentSchedule = paymentSchedule

	klogger.Exit(method)
	return nil
}

func (l *Loan) PerformPaymentCalc() error {
	method := "Loan.PerformPaymentCalc"
	klogger.Enter(method)

	err := l.ValidateCanPerformCalc()

	if err != nil {
		errMsg := "Loan payment calculation request is invalid:\n%v"
		klogger.ExitError(method, errMsg, err)
		return err
	}

	//int = (i+1)^n

	// Payment is P / {[(1+i)^n]-1} / [i(1+i)^n] where P is starting principal, i is the interest rate divided by 12, and n is the number of payments
	i := (l.InterestRate / 100) / 12
	p := l.Total
	n := l.LoanTerm
	l.MonthlyPayment = p / ((math.Pow((i+1), float64(n)) - 1) / (i * math.Pow((i+1), float64(n))))

	klogger.Exit(method)
	return nil
}

func (l *Loan) ValidateCanPerformCalc() error {
	method := "Loan.ValidateCanPerformCalc"
	klogger.Enter(method)

	if l.Total <= 0 {
		errMsg := "cannot perform calculation without total loan amount"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	if l.InterestRate < 0 {
		errMsg := "cannot perform calculation without interest rate"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	if l.LoanTerm <= 0 {
		errMsg := "cannot perform calculation without loan term"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	klogger.Exit(method)
	return nil
}

func (l *Loan) ValidateCanSaveLoan() error {
	method := "Loan.ValidateCanPerformCalc"
	klogger.Enter(method)

	if l.Name == "" {
		errMsg := "cannot save loan without loan name"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	if l.Total <= 0 {
		errMsg := "cannot save loan without total loan amount"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	if l.InterestRate < 0 {
		errMsg := "cannot save loan without interest rate"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	if l.LoanTerm <= 0 {
		errMsg := "cannot save loan without loan term"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	if l.UserID <= 0 {
		errMsg := "cannot save loan without userId"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	if l.TotalCost <= 0 {
		errMsg := "cannot save loan without total_cost"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	if l.TotalPayment <= 0 {
		errMsg := "cannot save loan without total payment"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	if l.Interest < 0 {
		errMsg := "cannot save loan without total interest"
		klogger.ExitError(method, errMsg)
		return errors.New(errMsg)
	}

	klogger.Exit(method)
	return nil
}

func (l *Loan) CompareLoanPayments(c Loan) []PaymentScheduleComparisonItem {
	method := "Loan.CompareLoanPayments"
	klogger.Enter(method)

	//m := int(math.Max(float64(len(l.PaymentSchedule)), float64(len(c.PaymentSchedule))))
	n := len(l.PaymentSchedule)
	m := len(c.PaymentSchedule)

	var r []PaymentScheduleComparisonItem
	var lPrincipalToDate float64
	var cPrincipalToDate float64
	var lInterestToDate float64
	var cInterestToDate float64

	for i := 0; i < n || i < m; i++ {

		var j PaymentScheduleComparisonItem

		if i < n && i < m {
			//Compare both elements
			lv := l.PaymentSchedule[i]
			cv := c.PaymentSchedule[i]

			j = PaymentScheduleComparisonItem{
				Month:                 lv.Month,
				Principal:             lv.Principal,
				PrincipalNew:          cv.Principal,
				PrincipalDelta:        cv.Principal - lv.Principal,
				Interest:              lv.Interest,
				InterestNew:           cv.Interest,
				InterestDelta:         cv.Interest - lv.Interest,
				InterestToDate:        lv.InterestToDate,
				InterestToDateNew:     cv.InterestToDate,
				InterestToDateDelta:   cv.InterestToDate - lv.InterestToDate,
				PrincipalToDate:       lv.PrincipalToDate,
				PrincipalToDateNew:    cv.PrincipalToDate,
				PrincipalToDateDelta:  cv.PrincipalToDate - lv.PrincipalToDate,
				RemainingBalance:      lv.RemainingBalance,
				RemainingBalanceNew:   cv.RemainingBalance,
				RemainingBalanceDelta: cv.RemainingBalance - lv.RemainingBalance,
			}

			lPrincipalToDate = lv.PrincipalToDate
			cPrincipalToDate = cv.PrincipalToDate
			lInterestToDate = lv.InterestToDate
			cInterestToDate = cv.InterestToDate

		} else if i < n {
			//Only original slice has entry
			lv := l.PaymentSchedule[i]

			j = PaymentScheduleComparisonItem{
				Month:                 lv.Month,
				Principal:             lv.Principal,
				PrincipalNew:          0,
				PrincipalDelta:        0 - lv.Principal,
				Interest:              lv.Interest,
				InterestNew:           0,
				InterestDelta:         0 - lv.Interest,
				InterestToDate:        lv.InterestToDate,
				InterestToDateNew:     cInterestToDate,
				InterestToDateDelta:   cInterestToDate - lv.InterestToDate,
				PrincipalToDate:       lv.PrincipalToDate,
				PrincipalToDateNew:    cPrincipalToDate,
				PrincipalToDateDelta:  cPrincipalToDate - lv.PrincipalToDate,
				RemainingBalance:      lv.RemainingBalance,
				RemainingBalanceNew:   0,
				RemainingBalanceDelta: 0 - lv.RemainingBalance,
			}

		} else {
			//Only new slice has entry
			cv := c.PaymentSchedule[i]

			j = PaymentScheduleComparisonItem{
				Month:                 cv.Month,
				Principal:             0,
				PrincipalNew:          cv.Principal,
				PrincipalDelta:        cv.Principal,
				Interest:              0,
				InterestNew:           cv.Interest,
				InterestDelta:         cv.Interest,
				InterestToDate:        lInterestToDate,
				InterestToDateNew:     cv.InterestToDate,
				InterestToDateDelta:   cv.InterestToDate - lInterestToDate,
				PrincipalToDate:       lPrincipalToDate,
				PrincipalToDateNew:    cv.PrincipalToDate,
				PrincipalToDateDelta:  cv.PrincipalToDate - lPrincipalToDate,
				RemainingBalance:      0,
				RemainingBalanceNew:   cv.RemainingBalance,
				RemainingBalanceDelta: cv.RemainingBalance,
			}
		}

		r = append(r, j)
	}

	klogger.Exit(method)
	return r
}
