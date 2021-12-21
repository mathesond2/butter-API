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
		fmt.Println(err, "zzz")
		u.Respond(w, u.Message(false, "GetTxn: Error while decoding request body"))
		return
	}

	fmt.Print(string(txn.WatchedAddress))

	blah := url.Values{
		"asset":          {txn.Asset},
		"status":         {txn.Status},
		"from":           {txn.From},
		"to":             {txn.To},
		"watchedAddress": {txn.WatchedAddress},
		"value":          {txn.Value},
		"direction":      {txn.Direction},
	}

	data, err := http.PostForm("http://stormy-cove-04196.herokuapp.com/api/associatedTxn", blah)
	if err != nil {
		fmt.Println("Error is req: ", err)
	}

	defer data.Body.Close()
	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		print("read body error: ", err)
	}

	fmt.Print(string(body))
	fmt.Println("data	", data)

	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
