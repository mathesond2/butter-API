package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-invoices/models"
	u "go-invoices/utils"
	"io/ioutil"
	"net/http"
	"strconv"
)

var ParseMempoolEvent = func(w http.ResponseWriter, r *http.Request) {
	txn := &models.Transaction{}

	err := json.NewDecoder(r.Body).Decode(txn)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "ParseMempoolEvent: Error while decoding request body"))
		return
	}

	parseValueStr, parseFloatErr := strconv.ParseFloat(txn.Value, 64)
	if parseFloatErr != nil {
		fmt.Println("parseFloatErr: ", parseFloatErr)
	}
	parsedTxnValue := parseValueStr / 1000000000000000000 //todo...do better than this

	params := make(map[string]interface{})
	params["asset"] = txn.Asset
	params["status"] = txn.Status
	params["from"] = txn.From
	params["to"] = txn.To
	params["watchedAddress"] = txn.WatchedAddress
	params["value"] = parsedTxnValue
	params["direction"] = txn.Direction

	bytesData, _ := json.Marshal(params)
	reader := bytes.NewReader(bytesData)

	request, error := http.NewRequest(
		http.MethodPost,
		"http://stormy-cove-04196.herokuapp.com/api/webhooks/associatedTxn",
		// "http://localhost:8000/api/webhooks/associatedTxn",
		reader,
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

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		print("error: ", err)
	}

	var data map[string]interface{}
	json.Unmarshal(body, &data)
	fmt.Printf("Results: %v\n", data)

	resp := u.Message(true, "success")
	resp["data"] = data["data"]
	u.Respond(w, resp)
}

var GetAssociatedTxn = func(w http.ResponseWriter, r *http.Request) {
	latestTxn := &models.ParsedTransaction{}

	err := json.NewDecoder(r.Body).Decode(latestTxn)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "GetAssociatedTxn: Error while decoding request body"))
		return
	}

	data := models.GetAssociatedTxn(latestTxn)

	// data.UserId
	GetWebhookByUserId(data.UserId)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

func GetWebhookByUserId(u uint) *models.Webhook {
	acc := &models.Webhook{}
	models.GetDB().Table("webhoooks").Where("user_id = ?", u).First(acc)
	if acc.UserId <= 0 { // no associated webhook found!
		return nil
	}

	return acc
}

// func GetWebhookByUserId(u uint) *models.Webhook {
// 	acc := &models.Webhook{}
// 	models.GetDB().Table("webhoooks").Where("id = ?", u).First(acc)
// 	if acc.UserId <= 0 { // no associated webhook found!
// 		return nil
// 	}

// 	return acc
// }
