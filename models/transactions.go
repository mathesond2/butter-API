package models

import (
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	Asset          string `json:"asset"`
	Status         string `json:"status"`
	From           string `json:"from"`
	To             string `json:"to"`
	WatchedAddress string `json:"watchedAddress"`
	Value          string `json:"value"`
	Direction      string `json:"direction"`
}

type ParsedTransaction struct {
	gorm.Model
	Asset          string  `json:"asset"`
	Status         string  `json:"status"`
	From           string  `json:"from"`
	To             string  `json:"to"`
	WatchedAddress string  `json:"watchedAddress"`
	Value          float64 `json:"value"`
	Direction      string  `json:"direction"`
}
