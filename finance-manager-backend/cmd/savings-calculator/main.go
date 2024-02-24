package main

import (
	"finance-manager-backend/internal/finance-mngr/enums/payfrequency"
	"finance-manager-backend/internal/finance-mngr/models/restmodels"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jon-kamis/klogger"
)

func main() {
	method := "savings-calculator.main"
	klogger.Enter(method)

	var r restmodels.SavingsCalculationRequest

	cmd := flag.NewFlagSet("calcSavings", flag.ExitOnError)
	goal := cmd.Float64("goal", 0.0, "savings goal")
	amount := cmd.Float64("amount", 0.0, "amount per pay to save")
	payFreq := cmd.String("payFreq", "", "frequency of pay. Options are 'weekly', 'bi-weekly', and 'monthly'")
	nPay := cmd.String("nextPay", "", "next payment date in yyyy-mm-dd format")
	d := cmd.String("deadline", "", "deadline for goal in yyyy-mm-dd format")

	cmd.Parse(os.Args[2:])

	if *nPay == "" || *payFreq == "" || *d == "" {
		panic("nextPay, payFrequency, and deadline are required")
	}

	r.Goal = *goal
	r.Amount = *amount
	r.PayFrequency = payfrequency.GetPayFrequency(*payFreq)
	r.NextPay = getDate(*nPay)
	r.Deadline = getDate(*d)

	resp, err := r.Calculate()

	if err != nil {
		panic(err)
	}

	klogger.Info(method, "Deadline: %v", resp.Deadline)
	klogger.Info(method, "Number of pays before deadline: %d", resp.NumPays)
	if r.Goal > 0.0 {
		klogger.Info(method, "Goal Details:")
		klogger.Info(method, "Goal: $%.2f", resp.Goal)
		klogger.Info(method, "Save per pay: $%.2f", resp.PerPay)
	}

	if resp.Actual > 0.0 {
		klogger.Info(method, "Savings Details:")
		klogger.Info(method, "Amount saving per pay: $%.2f", resp.ManPerPay)
		klogger.Info(method, "Amount saved by deadline: $%.2f", resp.Actual)
	}
	
	klogger.Exit(method)
}

func getDate(ds string) time.Time {
	method := "savings-calculator.getDate"
	klogger.Enter(method)

	dArr := strings.Split(ds, "-")

	y, err1 := strconv.Atoi(dArr[0])
	m, err2 := strconv.Atoi(dArr[1])
	d, err3 := strconv.Atoi(dArr[2])

	if err1 != nil {
		panic(err1)
	}

	if err2 != nil {
		panic(err2)
	} 

	if err3 != nil {
		panic(err3)
	}

	dt := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)

	klogger.Exit(method)
	return dt
}
