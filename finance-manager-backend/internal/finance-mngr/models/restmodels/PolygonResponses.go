package restmodels

type AggResponseItem struct {
	Ticker              string  `json:"T"`
	Close               float64 `json:"c"`
	High                float64 `json:"h"`
	Low                 float64 `json:"l"`
	Open                float64 `json:"o"`
	UnixTime            int     `json:"t"`
	TradeVolume         float64 `json:"v"`
	VolumeWeightedPrice float64 `json:"vw"`
}

type AggResponse struct {
	Adjusted     bool              `json:"adjusted"`
	QueryCount   int               `json:"queryCount"`
	RequestId    string            `json:"request_id"`
	Results      []AggResponseItem `json:"results"`
	ResultsCount int               `json:"resultsCount"`
	Status       string            `json:"status"`
	Ticker       string            `json:"ticker"`
}
