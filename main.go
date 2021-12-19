package main

import (
	"fmt"
	"go-contacts/app"
	"go-contacts/controllers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	router.HandleFunc("/api/invoices/new", controllers.CreateInvoice).Methods("POST")
	router.HandleFunc("/api/me/invoices", controllers.GetInvoices).Methods("GET")

	router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/api/{id}/contact", controllers.UpdateContact).Methods("PUT")
	router.HandleFunc("/api/{id}/contact", controllers.DeleteContact).Methods("DELETE")
	router.HandleFunc("/api/{id}/contact", controllers.GetContact).Methods("GET")
	router.HandleFunc("/api/me/contacts", controllers.GetContacts).Methods("GET") //  user/2/contacts

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
