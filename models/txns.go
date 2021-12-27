package models

import (
	"fmt"

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

func UpdateInvoiceStatusFromEvent(txn *ParsedTransaction) *Invoice {
	invoice := &Invoice{}

	//TODO: check invoice status for "in progress"
	err := GetDB().Table("invoices").Where(&Invoice{
		Sender_Address:    txn.WatchedAddress,
		Recipient_Address: txn.To,
		Amount:            txn.Value,
	}).First(&invoice).Error
	if err != nil {
		fmt.Println(err, "UpdateInvoiceStatusFromEvent")
		return nil
	}

	invoice.Status = txn.Status

	updatedErr := GetDB().Table("invoices").Where("id = ?", invoice.ID).Save(invoice).Error
	if updatedErr != nil {
		fmt.Println(updatedErr)
		return nil
	}

	return invoice
}
