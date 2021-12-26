package controllers

import (
	"encoding/json"
	"fmt"
	"go-invoices/models"
	u "go-invoices/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func PassesAddressChecks(invoice *models.Invoice, user uint) (bool, string) {
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

	return true, ""
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

	ok, msg := PassesAddressChecks(invoice, user)
	if !ok {
		u.Respond(w, u.Message(false, msg))
		return
	}

	resp := invoice.CreateInvoice()
	u.Respond(w, resp)
}

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	//should check if user is the owner of the invoice (user id)

	user := r.Context().Value("user").(uint)
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	invoice := &models.Invoice{}
	invoice.UserId = user

	err := json.NewDecoder(r.Body).Decode(invoice)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	resp := models.UpdateInvoice(id, invoice)
	u.Respond(w, resp)
}

func DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	//should check if user is the owner of the invoice (user id)

	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	data := models.DeleteInvoice(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

func GetInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	data := models.GetInvoice(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

func GetInvoices(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user").(uint64)
	data := models.GetInvoices(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
