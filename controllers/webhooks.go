package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-invoices/models"
	u "go-invoices/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type PostBody struct {
	ApiKey     string   `json:"apiKey"`
	Address    string   `json:"address"`
	Blockchain string   `json:"blockchain"`
	Networks   []string `json:"networks"`
}

type DeleteBody struct {
	Name string `json:"name"`
}

func AddressIsRegistered(address string, u uint) bool {
	resp := models.FindAddress(address, u)
	if resp["data"] == nil {
		return false
	} else {
		return true
	}
}

func WebhookIsRegistered(name string, u uint) bool {
	resp := models.FindWebhook(name, u)
	if resp["data"] == nil {
		return false
	} else {
		return true
	}
}

func WatchAddress(address string) string {
	supportedNetworks := []string{
		"main",
		"rinkeby",
	}

	postBody := PostBody{
		os.Getenv("blocknative_api_key"),
		address,
		"ethereum",
		supportedNetworks,
	}

	jsonData, jsonErr := json.Marshal(postBody)
	if jsonErr != nil {
		fmt.Println("error: ", jsonErr)
	}

	httpResp, httpErr := http.Post(
		"https://api.blocknative.com/address",
		"application/json; charset=utf-8",
		bytes.NewBuffer(jsonData),
	)
	if httpErr != nil {
		fmt.Println("watch address Error: ", httpErr)
	}

	defer httpResp.Body.Close()
	body, ioErr := ioutil.ReadAll(httpResp.Body)

	if ioErr != nil {
		fmt.Println("ioErr: ", ioErr)
	}

	type WatchAddressResponse struct {
		Msg string `json:"msg"`
	}
	var res WatchAddressResponse
	json.Unmarshal(body, &res)
	return res.Msg
}

//add address to the mempool watch list and then add to db
func AddAddress(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	address := &models.Address{}
	address.UserId = user

	err := json.NewDecoder(r.Body).Decode(address)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(false, "addAddress: Error while decoding request body"))
		return
	}

	if !u.IsValidEthAddress(address.Address) {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(false, "only valid Ethereum addresses are currently accepted"))
		return
	}

	if AddressIsRegistered(address.Address, user) {
		msg := fmt.Sprintf("address %s is already registered with your account.", address.Address)
		u.Respond(w, u.Message(false, msg))
		return
	}

	watchAddressRes := WatchAddress(address.Address)
	if watchAddressRes != "success" { //response body from blocknative endpoint
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(false, watchAddressRes))
		return
	}

	resp := models.CreateAddress(address)
	u.Respond(w, resp)
}

func GetAddresses(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user").(uint)
	data := models.GetAddresses(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

func UnwatchAddress(address string) string {
	supportedNetworks := []string{
		"main",
		"rinkeby",
	}

	params := make(map[string]interface{})
	params["apiKey"] = os.Getenv("blocknative_api_key")
	params["address"] = address
	params["blockchain"] = "ethereum"
	params["networks"] = supportedNetworks

	bytesData, _ := json.Marshal(params)
	reader := bytes.NewReader(bytesData)

	httpReq, httpErr := http.NewRequest(
		http.MethodDelete,
		"https://api.blocknative.com/address",
		reader,
	)
	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if httpErr != nil {
		fmt.Println("httpReq Error: ", httpErr)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if httpErr != nil {
		fmt.Println("watch address Error: ", httpErr)
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, ioErr := ioutil.ReadAll(resp.Body)
	if ioErr != nil {
		fmt.Println("ioErr: ", ioErr)
	}

	type WatchAddressResponse struct {
		Msg string `json:"msg"`
	}
	var res WatchAddressResponse
	json.Unmarshal(body, &res)
	return res.Msg
}

//delete address from watch list and then from db
func DeleteAddress(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	address := &models.Address{}
	address.UserId = user

	err := json.NewDecoder(r.Body).Decode(address)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "DeleteAddress: Error while decoding request body"))
		return
	}

	if !u.IsValidEthAddress(address.Address) {
		u.Respond(w, u.Message(false, "only valid Ethereum addresses are currently accepted"))
		return
	}

	if !AddressIsRegistered(address.Address, user) {
		msg := fmt.Sprintf("address %s is not registered with your account.", address.Address)
		u.Respond(w, u.Message(false, msg))
		return
	}

	unwatchAddressRes := UnwatchAddress(address.Address)
	if unwatchAddressRes != "success" {
		u.Respond(w, u.Message(false, unwatchAddressRes))
		return
	}

	resp := models.DeleteAddress(address)
	u.Respond(w, resp)
}

func AddWebhook(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	webhook := &models.Webhook{}
	webhook.UserId = user

	err := json.NewDecoder(r.Body).Decode(webhook)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		u.Respond(w, u.Message(false, "AddWebhook: Error while decoding request body"))
		return
	}

	if !AddressIsRegistered(webhook.Address, user) {
		msg := fmt.Sprintf("address %s is not registered with your account.", webhook.Address)
		u.Respond(w, u.Message(false, msg))
		return
	}

	isRegisteredWebhook := WebhookIsRegistered(webhook.Name, user)
	if isRegisteredWebhook {
		msg := fmt.Sprintf("webhook name '%s' is already registered with your account, please choose another name", webhook.Name)
		u.Respond(w, u.Message(false, msg))
		return
	}

	resp := models.CreateWebhook(webhook)
	if resp["status"] == false {
		w.WriteHeader(http.StatusNotFound)
	}

	u.Respond(w, resp)
}

func DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	webhookReq := &DeleteBody{}

	err := json.NewDecoder(r.Body).Decode(webhookReq)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		u.Respond(w, u.Message(false, "DeleteWebhook: Error while decoding request body"))
		return
	}

	if webhookReq.Name == "" {
		w.WriteHeader(http.StatusNotFound)
		u.Respond(w, u.Message(false, "name property is required"))
		return
	}

	resp := models.DeleteWebhook(webhookReq.Name, user)
	u.Respond(w, resp)
}

func GetWebhooks(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user").(uint)
	data := models.GetWebhooks(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

type WebhookReqBody struct {
	Name     string          `json:"name"`
	Networks string          `json:"networks"`
	Invoice  *models.Invoice `json:"invoice"`
}

func SendDataToWebhook(w models.Webhook, invoice *models.Invoice) {
	postBody := WebhookReqBody{
		w.Name,
		w.Networks,
		invoice,
	}

	jsonData, jsonErr := json.Marshal(postBody)
	if jsonErr != nil {
		fmt.Println("error: ", jsonErr)
	}

	httpResp, httpErr := http.Post(
		w.Endpoint_Url,
		"application/json; charset=utf-8",
		bytes.NewBuffer(jsonData),
	)

	if httpErr != nil {
		fmt.Println("watch address Error: ", httpErr)
	}

	defer httpResp.Body.Close()
}
