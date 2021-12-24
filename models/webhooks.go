package models

import (
	u "go-invoices/utils"
	"strings"

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

func (w *PreParsedWebhook) ValidateWebhook() (map[string]interface{}, bool) {
	if w.Address == "" {
		return u.Message(false, "address should be on the payload"), false
	}

	if len(w.Networks) == 0 {
		return u.Message(false, "at least one chosen network should be on the payload"), false
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

func CreateWebhook(w *PreParsedWebhook) map[string]interface{} {
	if resp, ok := w.ValidateWebhook(); !ok {
		return resp
	}

	hackyArrToStr := strings.Join(w.Networks, " ")

	parsedWebhook := &Webhook{
		Address:      w.Address,
		Networks:     hackyArrToStr,
		Name:         w.Name,
		Endpoint_Url: w.Endpoint_Url,
		UserId:       w.UserId,
	}

	GetDB().Create(parsedWebhook)

	resp := u.Message(true, "success")
	resp["data"] = w
	return resp
}
