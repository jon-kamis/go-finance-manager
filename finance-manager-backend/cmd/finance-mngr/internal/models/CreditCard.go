package models

import "time"

type CreditCard struct {
	ID                   int       `json:"id"`
	UserID               int       `json:"userId" gorm:"user_id"`
	Name                 string    `json:"name"`
	Balance              float64   `json:"balance"`
	APR                  float64   `json:"apr"`
	MinPayment           float64   `json:"minPayment" gorm:"column:min_pay" `
	MinPaymentPercentage float64   `json:"minPaymentPercentage" gorm:"column:min_pay_percentage"`
	CreateDt             time.Time `json:"createDt" gorm:"column:create_dt"`
	LastUpdateDt         time.Time `json:"lastUpdateDt" gorm:"column:last_update_dt"`
}
