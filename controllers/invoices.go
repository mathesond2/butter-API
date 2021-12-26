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

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	invoice := &models.Invoice{}

	err := json.NewDecoder(r.Body).Decode(invoice)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	invoice.UserId = user
	resp := invoice.CreateInvoice()
	u.Respond(w, resp)
}

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	//should check if user is the owner of the invoice (user id)

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
