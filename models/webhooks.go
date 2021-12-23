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

func (w *Webhook) ValidateWebhook() (map[string]interface{}, bool) {
	if w.Address == "" {
		return u.Message(false, "address should be on the payload"), false
	}

	if w.Networks == "" {
		return u.Message(false, "chosen networks should be on the payload"), false
	}

	if w.Name == "" {
		return u.Message(false, "webhook name should be on the payload"), false
	}

	if w.Endpoint_Url == "" {
		return u.Message(false, "endpoint_url should be on the payload"), false
	}

	if w.UserId <= 0 {
		return u.Message(false, "User is not recognized"), false
	}

	return u.Message(true, "success"), true
}

func CreateWebhook(w *Webhook) map[string]interface{} {
	if resp, ok := w.ValidateWebhook(); !ok {
		return resp
	}

	GetDB().Create(w)

	resp := u.Message(true, "success")
	resp["data"] = w
	return resp
}
