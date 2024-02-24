package payfrequency

import "finance-manager-backend/internal/finance-mngr/constants"

type PayFrequency string

const (
	Undefined PayFrequency = ""
	Weekly    PayFrequency = constants.IncomeFreqWeekly
	BiWeekly  PayFrequency = constants.IncomeFreqBiWeekly
	Monthly   PayFrequency = constants.IncomeFreqMonthly
)

func GetPayFrequency(s string) PayFrequency {
	switch s {
	case constants.IncomeFreqWeekly:
		return Weekly
	case constants.IncomeFreqBiWeekly:
		return BiWeekly
	case constants.IncomeFreqMonthly:
		return Monthly
	default:
		return Undefined
	}
}

func (p PayFrequency) String() string {
	switch p {
	case Weekly:
		return constants.IncomeFreqWeekly
	case BiWeekly:
		return constants.IncomeFreqBiWeekly
	case Monthly:
		return constants.IncomeFreqMonthly
	default:
		return ""
	}
}
