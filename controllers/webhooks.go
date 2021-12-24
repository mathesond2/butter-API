package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-invoices/models"
	u "go-invoices/utils"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type PostBody struct {
	ApiKey     string   `json:"apiKey"`
	Address    string   `json:"address"`
	Blockchain string   `json:"blockchain"`
	Networks   []string `json:"networks"`
}

//this allows us to add an address to the mempool watch list
func AddAddressToWatch(w http.ResponseWriter, r *http.Request) {
	addressAuth := &models.AddressAuth{}

	err := json.NewDecoder(r.Body).Decode(addressAuth)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	if len(addressAuth.Addresses) == 0 {
		u.Respond(w, u.Message(false, "Invalid request: no addresses provided"))
		return
	}

	//we should also look for any dupes in the db
	for _, address := range addressAuth.Addresses {
		if !strings.HasPrefix(address, "0x") {
			u.Respond(w, u.Message(false, "only valid Ethereum addresses are currently accepted"))
			return
		}
	}

	var addressReqResults = make(map[string]string)

	for _, address := range addressAuth.Addresses {
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

		resp, httpErr := http.Post(
			"https://api.blocknative.com/address",
			"application/json; charset=utf-8",
			bytes.NewBuffer(jsonData),
		)

		if httpErr != nil {
			fmt.Println("watch address Error: ", httpErr)
		}

		defer resp.Body.Close()
		body, ioErr := ioutil.ReadAll(resp.Body)
		if ioErr != nil {
			fmt.Println("ioErr: ", ioErr)
		}

		type Response struct {
			Msg string `json:"msg"`
		}
		var res Response
		json.Unmarshal(body, &res)
		addressReqResults[address] = res.Msg
	}

	resp := u.Message(true, "success")
	resp["data"] = addressReqResults
	u.Respond(w, resp)
}

func AddWebhook(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	webhook := &models.PreParsedWebhook{}

	err := json.NewDecoder(r.Body).Decode(webhook)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "addWebhook: Error while decoding request body"))
		return
	}

	webhook.UserId = user

	resp := models.CreateWebhook(webhook)
	u.Respond(w, resp)
}

func GetWebhookByUserId(u uint) *models.Webhook {
	acc := &models.Webhook{}

	err := models.GetDB().Table("webhooks").Where("user_id = ?", u).First(acc).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("webhook: ", acc)
	return acc
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
