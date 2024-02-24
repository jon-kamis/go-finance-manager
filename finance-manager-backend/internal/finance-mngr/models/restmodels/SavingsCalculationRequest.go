package restmodels

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/enums/payfrequency"
	"math"
	"time"

	"github.com/jon-kamis/klogger"
)

type SavingsCalculationRequest struct {
	Goal         float64                   `json:"goal"`
	Amount       float64                   `json:"amount"`
	Deadline     time.Time                 `json:"deadline"`
	PayFrequency payfrequency.PayFrequency `json:"payFrequency"`
	NextPay      time.Time                 `json:"nextPay"`
}

func (s *SavingsCalculationRequest) Calculate() (SavingsCalculationResponse, error) {
	method := "SavingsCalculationRequest.Calculate"
	klogger.Enter(method)

	var r SavingsCalculationResponse

	if s.Goal == 0 && s.Amount == 0 {
		err := errors.New("goal or amount is required")
		klogger.ExitError(method, err.Error())
		return r, err
	}

	if s.Deadline.IsZero() {
		err := errors.New("deadline is required")
		klogger.ExitError(method, err.Error())
		return r, err
	}

	if s.PayFrequency == payfrequency.Undefined {
		err := errors.New("pay frequency is required")
		klogger.ExitError(method, err.Error())
		return r, err
	}

	r.Deadline = s.Deadline
	r.ManPerPay = s.Amount

	//Get number of pays
	r.NumPays = s.GetPaysBeforeDeadline()
	pays := float64(r.NumPays)

	if s.Goal != 0.0 {
		r.Goal = s.Goal
		r.PerPay = (math.Round((s.Goal / pays) * 100)) / 100
	}

	if s.Amount > 0 {
		r.Actual = pays * s.Amount
	}

	klogger.Exit(method)
	return r, nil
}

func (r *SavingsCalculationRequest) GetPaysBeforeDeadline() int {
	method := "SavingsCalculationRequest.GetPaysBeforeDeadline"
	klogger.Enter(method)

	if r.Deadline.Before(time.Now()) {
		klogger.Info(method, "deadline is a past date")
		klogger.Exit(method)
		return 0
	}

	dt := r.NextPay

	if dt.After(r.Deadline) {
		klogger.Info(method, "next pay is after deadline")
		klogger.Exit(method)
		return 0
	}

	numPays := 0

	if r.PayFrequency == payfrequency.Monthly {

		for !dt.After(r.Deadline) {
			numPays++
			dt = dt.AddDate(0, 1, 0)
		}

		klogger.Exit(method)
		return numPays
	} else if r.PayFrequency == payfrequency.Weekly {
		for !dt.After(r.Deadline) {
			numPays++
			dt = dt.AddDate(0, 0, 7)
		}

		klogger.Exit(method)
		return numPays
	} else {
		for !dt.After(r.Deadline) {
			numPays++
			dt = dt.AddDate(0, 0, 14)
		}

		klogger.Exit(method)
		return numPays
	}
}
