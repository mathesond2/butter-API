package main

import (
	"fmt"
	"go-invoices/app"
	"go-invoices/controllers"
	m "go-invoices/models"
	u "go-invoices/utils"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var NotFound = func(w http.ResponseWriter, r *http.Request) {
	u.Respond(w, u.Message(false, "Invalid request: route not found"))
}

func main() {
	m.GetDB()

	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	router.HandleFunc("/api/webhooks/add", controllers.AddWebhook).Methods("POST")
	router.HandleFunc("/api/webhooks/address/add", controllers.AddAddressToWatch).Methods("POST")
	router.HandleFunc("/api/webhooks/mempoolEvent", controllers.ParseMempoolEvent).Methods("POST")
	router.HandleFunc("/api/webhooks/updateInvoice", controllers.UpdateInvoiceFromEvent).Methods("POST") //prob should be put

	router.HandleFunc("/api/invoice/new", controllers.CreateInvoice).Methods("POST")
	router.HandleFunc("/api/{id}/invoice", controllers.UpdateInvoice).Methods("PUT")
	router.HandleFunc("/api/{id}/invoice", controllers.DeleteInvoice).Methods("DELETE")
	router.HandleFunc("/api/{id}/invoice", controllers.GetInvoice).Methods("GET")

	router.HandleFunc("/api/me/invoices", controllers.GetInvoices).Methods("GET")

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	router.NotFoundHandler = http.HandlerFunc(NotFound)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
