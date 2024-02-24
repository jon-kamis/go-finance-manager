package restmodels

import "time"

type SavingsCalculationResponse struct {
	Goal      float64   `json:"goal"`
	PerPay    float64   `json:"perPay"`
	ManPerPay float64   `json:"manPerPay"`
	Actual    float64   `json:"actual"`
	NumPays   int       `json:"numPays"`
	Deadline  time.Time `json:"deadline"`
}
