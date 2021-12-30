package controllers

import (
	"encoding/json"
	"fmt"
	"go-invoices/models"
	u "go-invoices/utils"
	"net/http"
	"strconv"
)

func hasRegisteredAddress(a1 string, a2 string, user uint) bool {
	res1 := AddressIsRegistered(a1, user)
	res2 := AddressIsRegistered(a2, user)
	if res1 || res2 {
		return true
	} else {
		return false
	}
}

func PassesAddressAndWebhookChecks(invoice *models.Invoice, user uint) (bool, string) {
	if len(invoice.Recipient_Address) == 0 || len(invoice.Sender_Address) == 0 {
		return false, "both sender and recipient addresses must be on the payload"
	}

	if !u.IsValidEthAddress(invoice.Recipient_Address) || !u.IsValidEthAddress(invoice.Sender_Address) {
		return false, "only valid Ethereum addresses are currently accepted"
	}

	isRegistered := hasRegisteredAddress(
		invoice.Recipient_Address,
		invoice.Sender_Address,
		user,
	)

	if !isRegistered {
		return false, "At least one of the addresses provided must be registered with your account"
	}

	if len(invoice.Webhook_Name) > 0 {
		isRegisteredWebhook := WebhookIsRegistered(invoice.Webhook_Name, user)
		if !isRegisteredWebhook {
			return false, "The webhook name provided is not registered with your account"
		}
	}

	return true, ""
}

func UpdateInvoiceStatusAsPending(w http.ResponseWriter, r *http.Request) {
	userInfo := &models.UserInfo{}
	err := json.NewDecoder(r.Body).Decode(userInfo)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	data := models.UpdateInvoiceStatusAsPending(userInfo.ID)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	invoice := &models.Invoice{}
	invoice.UserId = user

	err := json.NewDecoder(r.Body).Decode(invoice)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	ok, msg := PassesAddressAndWebhookChecks(invoice, user)
	if !ok {
		u.Respond(w, u.Message(false, msg))
		return
	}

	resp := invoice.CreateInvoice()
	u.Respond(w, resp)
}

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	invoice := &models.Invoice{}
	invoice.UserId = user

	err := json.NewDecoder(r.Body).Decode(invoice)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	ok, msg := PassesAddressAndWebhookChecks(invoice, user)
	if !ok {
		u.Respond(w, u.Message(false, msg))
		return
	}

	resp := models.UpdateInvoice(invoice)
	u.Respond(w, resp)
}

func DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	userInfo := &models.UserInfo{}
	userInfo.UserId = user

	err := json.NewDecoder(r.Body).Decode(userInfo)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	data, err := models.DeleteInvoice(userInfo)

	var resp map[string]interface{}
	if err != nil {
		resp = u.Message(false, err.Error())
	} else {
		resp = u.Message(true, "success")
		resp["data"] = data
	}

	u.Respond(w, resp)
}

func GetInvoice(w http.ResponseWriter, r *http.Request) {
	userInfo := &models.UserInfo{}
	givenId := r.FormValue("id")

	if len(givenId) == 0 {
		w.WriteHeader(http.StatusNotFound)
		u.Respond(w, u.Message(false, "GetInvoice: id is required as query parameter"))
		return
	}

	u64, err := strconv.ParseUint(givenId, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		u.Respond(w, u.Message(false, "GetInvoice: query parameter 'id' must be a number"))
		return
	}

	userInfo.ID = uint(u64)

	data := models.GetInvoice(userInfo.ID)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

func GetInvoices(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user").(uint)
	data := models.GetInvoices(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
