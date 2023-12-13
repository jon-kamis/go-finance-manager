package models

type Bill struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Name   string `json:"name"`
	Amount string `json:"amount"`
}
