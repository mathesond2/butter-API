package controllers

import (
	"encoding/json"
	"fmt"
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var CreateContact = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	contact := &models.Contact{}

	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	contact.UserId = user
	resp := contact.Create()
	u.Respond(w, resp)
}

var GetContactsFor = func(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("user").(uint)
	data := models.GetContacts(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var GetContact = func(w http.ResponseWriter, r *http.Request) {
	// str_uint := "1234"
	// val_uint, _ := strconv.ParseUint(str_uint, 10, 64)
	// fmt.Printf("%d\n", val_uint)

	vars := mux.Vars(r)
	id := vars["id"]
	idNum, _ := strconv.ParseUint(id, 10, 32)
	fmt.Println("zzz", idNum)

	// id := r.Context().Value("id").(uint)
	data := models.GetContact(idNum)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
