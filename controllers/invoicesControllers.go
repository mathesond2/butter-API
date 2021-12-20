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

var CreateInvoice = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	invoice := &models.Invoice{}

	err := json.NewDecoder(r.Body).Decode(invoice)
	if err != nil {
		fmt.Println(err, "zzz")
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	invoice.UserId = user
	resp := invoice.CreateInvoice()
	u.Respond(w, resp)
}

var UpdateInvoice = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)

	invoice := &models.Invoice{}

	err := json.NewDecoder(r.Body).Decode(invoice)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	invoice.UserId = user
	resp := models.UpdateInvoice(id, invoice)
	u.Respond(w, resp)
}

var DeleteInvoice = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	data := models.DeleteInvoice(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var GetInvoice = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	data := models.GetInvoice(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var GetInvoices = func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user").(uint)
	data := models.GetInvoices(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
