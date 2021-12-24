package main

import (
	"fmt"
	"go-invoices/app"
	c "go-invoices/controllers"
	m "go-invoices/models"
	u "go-invoices/utils"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var NotFound = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	u.Respond(w, u.Message(false, "This resources was not found on our server"))
}

func main() {
	m.GetDB()

	r := mux.NewRouter()

	r.HandleFunc("/api/user", c.CreateAccount).Methods("POST")
	r.HandleFunc("/api/user/login", c.Authenticate).Methods("POST")

	r.HandleFunc("/api/webhooks", c.AddWebhook).Methods("POST")
	r.HandleFunc("/api/webhooks/address", c.AddAddressToWatch).Methods("POST")
	r.HandleFunc("/api/webhooks/mempoolEvent", c.ParseMempoolEvent).Methods("POST")
	r.HandleFunc("/api/webhooks/updateInvoice", c.UpdateInvoiceFromEvent).Methods("POST") //prob should be put

	r.HandleFunc("/api/invoice", c.CreateInvoice).Methods("POST")
	r.HandleFunc("/api/{id}/invoice", c.UpdateInvoice).Methods("PUT")
	r.HandleFunc("/api/{id}/invoice", c.DeleteInvoice).Methods("DELETE")
	r.HandleFunc("/api/{id}/invoice", c.GetInvoice).Methods("GET")

	r.HandleFunc("/api/invoices", c.GetInvoices).Methods("GET")

	r.Use(app.JwtAuthentication)

	r.NotFoundHandler = http.HandlerFunc(NotFound)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Print(err)
	}
}
