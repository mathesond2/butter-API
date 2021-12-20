package controllers

import (
	"encoding/json"
	"fmt"
	"go-invoices/models"
	u "go-invoices/utils"
	"net/http"
)

var GetTxn = func(w http.ResponseWriter, r *http.Request) {
	txn := &models.Transaction{}

	err := json.NewDecoder(r.Body).Decode(txn)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	// invoice.UserId = user
	// resp := invoice.CreateInvoice()
	fmt.Println("GetTxn", txn)
	resp := u.Message(true, "success")
	u.Respond(w, resp)
}
