package constants

const IncomeTypeSalary = "salary"
const IncomeTypeHourly = "hourly"
const IncomeFreqWeekly = "weekly"
const IncomeFreqBiWeekly = "bi-weekly"
const IncomeFreqMonthly = "monthly"

var ValidTypes = []string{IncomeTypeHourly, IncomeTypeSalary}
var ValidFreq = []string{IncomeFreqWeekly, IncomeFreqBiWeekly, IncomeFreqMonthly}
