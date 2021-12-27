package models

import (
	"fmt"
	u "go-invoices/utils"

	"github.com/jinzhu/gorm"
)

type Webhook struct {
	gorm.Model
	Address      string `json:"address"`
	Networks     string `json:"networks"`
	Name         string `json:"name"`
	Endpoint_Url string `json:"endpoint_url"`
	UserId       uint   `json:"user_id"`
}

type Address struct {
	gorm.Model
	Address string `json:"address"`
	UserId  uint   `json:"user_id"`
}

type AddressAuth struct {
	gorm.Model
	Addresses []string `json:"addresses"`
}

func (w *Webhook) ValidateWebhook() (map[string]interface{}, bool) {
	if w.Address == "" {
		return u.Message(false, "address should be on the payload"), false
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

	// networks := [2]string{"main", "rinkeby"} //until we need a user to specify
	webhookWithNetworks := &Webhook{
		Address:      w.Address,
		Networks:     "main rinkeby",
		Name:         w.Name,
		Endpoint_Url: w.Endpoint_Url,
		UserId:       w.UserId,
	}

	GetDB().Create(webhookWithNetworks)

	resp := u.Message(true, "success")
	resp["data"] = w
	return resp
}

func CreateAddress(a *Address) map[string]interface{} {
	GetDB().Create(a)
	resp := u.Message(true, "success")
	resp["data"] = a
	return resp
}

func DeleteAddress(a *Address) map[string]interface{} {
	addy := &Address{}

	deleteErr := GetDB().Table("addresses").Where(&Address{
		Address: a.Address,
		UserId:  a.UserId,
	}).Delete(&addy).Error
	if deleteErr != nil {
		fmt.Println("DeleteInvoice delete error: ", deleteErr)
		return nil
	}

	resp := u.Message(true, "success")
	msg := fmt.Sprintf("address %s has been deleted.", a.Address)
	resp["data"] = msg
	return resp
}

func FindAddress(a string, user uint) map[string]interface{} {
	address := &Address{}

	err := GetDB().Table("addresses").Where(&Address{
		Address: a,
		UserId:  user,
	}).First(&address).Error
	if err != nil {
		return nil
	}

	resp := u.Message(true, "success")
	resp["data"] = address.Address
	return resp
}
