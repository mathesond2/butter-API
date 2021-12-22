package models

import (
	"fmt"
	u "go-invoices/utils"

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

func GetAssociatedTxn(txn *ParsedTransaction) *Invoice {
	invoice := &Invoice{}

	err := GetDB().Table("invoices").Where(&Invoice{
		Sender_Address:    txn.WatchedAddress,
		Recipient_Address: txn.To,
		Amount:            txn.Value,
	}).First(&invoice).Error
	if err != nil {
		fmt.Println(err, "GetAssociatedTxn")
		return nil
	}

	var status string
	if txn.Status == "confirmed" {
		status = "confirmed"
	} else {
		status = "pending"
	}
	invoice.Status = status

	updatedErr := GetDB().Table("invoices").Where("id = ?", invoice.ID).Save(invoice).Error
	if updatedErr != nil {
		fmt.Println(updatedErr)
		return nil
	}

	return invoice
}

type Webhook struct {
	Address     string   `json:"address"`
	Networks    []string `json:"networks"`
	Name        string   `json:"name"`
	EndpointUrl string   `json:"endpointUrl"`
	UserId      uint     `json:"userId"`
}

func (webhook *Webhook) ValidateWebhook() (map[string]interface{}, bool) {
	if webhook.Address == "" {
		return u.Message(false, "address should be on the payload"), false
	}

	if len(webhook.Networks) == 0 {
		return u.Message(false, "chosen networks should be on the payload"), false
	}

	if webhook.Name == "" {
		return u.Message(false, "webhook name should be on the payload"), false
	}

	if webhook.EndpointUrl == "" {
		return u.Message(false, "endpointUrl should be on the payload"), false
	}

	if webhook.UserId <= 0 {
		return u.Message(false, "User is not recognized"), false
	}

	return u.Message(true, "success"), true
}

func CreateWebhook(webhook *Webhook) map[string]interface{} {
	if resp, ok := webhook.ValidateWebhook(); !ok {
		return resp
	}

	GetDB().Create(webhook)

	resp := u.Message(true, "success")
	resp["data"] = webhook
	return resp
}
