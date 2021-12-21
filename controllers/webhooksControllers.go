package controllers

import (
	"encoding/json"
	"fmt"
	"go-invoices/models"
	u "go-invoices/utils"
	"io/ioutil"
	"net/http"
	"net/url"
)

var GetTxn = func(w http.ResponseWriter, r *http.Request) {
	txn := &models.Transaction{}

	err := json.NewDecoder(r.Body).Decode(txn)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	fmt.Print(string(txn.WatchedAddress))
	data, err := http.PostForm("http://stormy-cove-04196.herokuapp.com/api/associatedTxn", url.Values{
		"asset":          {txn.Asset},
		"status":         {txn.Status},
		"from":           {txn.From},
		"to":             {txn.To},
		"watchedAddress": {txn.WatchedAddress},
		"value":          {txn.Value},
		"direction":      {txn.Direction},
	})
	if err != nil {
		fmt.Println("Error is req: ", err)
	}

	defer data.Body.Close()
	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		print(err)
	}

	fmt.Print(string(body))
	fmt.Println("data	", data)

	resp := u.Message(true, "success")
	u.Respond(w, resp)
}
