package models

import (
	u "go-invoices/utils"

	"github.com/jinzhu/gorm"
)

type PreParsedWebhook struct {
	Address      string   `json:"address"`
	Networks     []string `json:"networks"`
	Name         string   `json:"name"`
	Endpoint_Url string   `json:"endpoint_url"`
	UserId       uint     `json:"user_id"`
}
type Webhook struct {
	gorm.Model
	Address      string `json:"address"`
	Networks     string `json:"networks"`
	Name         string `json:"name"`
	Endpoint_Url string `json:"endpoint_url"`
	UserId       uint   `json:"user_id"`
}

func (webhook *Webhook) ValidateWebhook() (map[string]interface{}, bool) {
	if webhook.Address == "" {
		return u.Message(false, "address should be on the payload"), false
	}

	if webhook.Networks == "" {
		return u.Message(false, "chosen networks should be on the payload"), false
	}

	if webhook.Name == "" {
		return u.Message(false, "webhook name should be on the payload"), false
	}

	if webhook.Endpoint_Url == "" {
		return u.Message(false, "endpoint_url should be on the payload"), false
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
