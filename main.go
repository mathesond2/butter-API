package main

import (
	"fmt"
	"go-invoices/app"
	"go-invoices/controllers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	router.HandleFunc("/api/invoice/new", controllers.CreateInvoice).Methods("POST")
	router.HandleFunc("/api/{id}/invoice", controllers.UpdateInvoice).Methods("PUT")
	router.HandleFunc("/api/{id}/invoice", controllers.DeleteInvoice).Methods("DELETE")
	router.HandleFunc("/api/{id}/invoice", controllers.GetInvoice).Methods("GET")
	router.HandleFunc("/api/me/invoices", controllers.GetInvoices).Methods("GET")

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
