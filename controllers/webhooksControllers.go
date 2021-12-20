package controllers

import (
	"fmt"
	u "go-invoices/utils"
	"net/http"
)

var GetTxn = func(w http.ResponseWriter, r *http.Request) {
	// invoice := &models.Invoice{}

	// err := json.NewDecoder(r.Body).Decode(invoice)
	// if err != nil {
	// 	u.Respond(w, u.Message(false, "Error while decoding request body"))
	// 	return
	// }

	// invoice.UserId = user
	// resp := invoice.CreateInvoice()
	resp := u.Message(true, "success")
	u.Respond(w, resp)
	fmt.Println("GetTxn", r.Body)
}
