package controllers

import (
	"bytes"
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
		fmt.Println(err)
		u.Respond(w, u.Message(false, "GetTxn: Error while decoding request body"))
		return
	}

	fmt.Print(string(txn.WatchedAddress))

	// blah := url.Values{
	// 	"asset":          {txn.Asset},
	// 	"status":         {txn.Status},
	// 	"from":           {txn.From},
	// 	"to":             {txn.To},
	// 	"watchedAddress": {txn.WatchedAddress},
	// 	"value":          {txn.Value},
	// 	"direction":      {txn.Direction},
	// }

	params := make(map[string]interface{})
	params["asset"] = txn.Asset
	params["status"] = txn.Status
	params["from"] = txn.From
	params["to"] = txn.To
	params["watchedAddress"] = txn.WatchedAddress
	params["value"] = txn.Value
	params["direction"] = txn.Direction

	bytesData, _ := json.Marshal(params)
	reader := bytes.NewReader(bytesData)

	// data, err := http.PostForm("http://stormy-cove-04196.herokuapp.com/api/associatedTxn", blah)
	request, error := http.NewRequest(
		http.MethodPost,
		"http://stormy-cove-04196.herokuapp.com/api/associatedTxn",
		reader,
		// strings.NewReader(blah.Encode()),
	)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if error != nil {
		fmt.Println("Error is req: ", error)
	}

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	invoice := &models.Invoice{}
	decodeInvoiceErr := json.NewDecoder(response.Body).Decode(invoice)
	if decodeInvoiceErr != nil {
		fmt.Println(decodeInvoiceErr)
		u.Respond(w, u.Message(false, "GetTxn: Error while decoding request body"))
		return
	}
	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	print("read body error: ", err)
	// }

	// fmt.Println("response Body:", string(body))

	// defer data.Body.Close()
	// body, err := ioutil.ReadAll(data.Body)
	// if err != nil {
	// 	print("read body error: ", err)
	// }

	// fmt.Print(string(body))
	// ioutil.ReadAll
	resp := u.Message(true, "success")
	resp["data"] = invoice
	u.Respond(w, resp)
}
