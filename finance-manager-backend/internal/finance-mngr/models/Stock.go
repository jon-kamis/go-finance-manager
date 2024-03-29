package models

import "time"

//Contains Data for a Stock object
type Stock struct {
	ID           int       `json:"id" gorm:"id"`
	Ticker       string    `json:"ticker"`
	High         float64   `json:"high" gorm:"column:high"`
	Low          float64   `json:"low" gorm:"column:low"`
	Open         float64   `json:"open" gorm:"column:open"`
	Close        float64   `json:"close" gorm:"column:close"`
	Date         time.Time `json:"date" gorm:"column:date"`
	CreateDt     time.Time `json:"createDt"`
	LastUpdateDt time.Time `json:"lastUpdateDt"`
}

type StockData struct {
	ID           int       `json:"id" gorm:"id"`
	Ticker       string    `json:"ticker"`
	High         float64   `json:"high" gorm:"column:high"`
	Low          float64   `json:"low" gorm:"column:low"`
	Open         float64   `json:"open" gorm:"column:open"`
	Close        float64   `json:"close" gorm:"column:close"`
	Date         time.Time `json:"-" gorm:"column:date"`
	CreateDt     time.Time `json:"createDt"`
	LastUpdateDt time.Time `json:"lastUpdateDt"`
}
