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

var AddWallet = func(w http.ResponseWriter, r *http.Request) {
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

	// var resp models.Response

	var addressReqResults = make(map[string]string)
	//this should be its own validate fn where we also look up any dupes in the db
	for _, address := range addressAuth.Addresses {
		if !strings.HasPrefix(address, "0x") {
			u.Respond(w, u.Message(false, "only valid Ethereum addresses are currently accepted"))
			return
		}
	}

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

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create()
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}
