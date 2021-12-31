package main

import (
	"fmt"
	"go-invoices/app"
	c "go-invoices/controllers"
	u "go-invoices/utils"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var NotFound = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	u.Respond(w, u.Message(false, "This resource was not found on our server"))
}

func main() {
	r := mux.NewRouter()

	//Private
	r.HandleFunc("/webhooks/mempoolEvent", c.ParseMempoolEvent).Methods("POST")
	r.HandleFunc("/webhooks/updateInvoiceStatusFromEvent", c.UpdateInvoiceStatusFromEvent).Methods("POST") //prob should be put
	r.HandleFunc("/invoice/status", c.UpdateInvoiceStatusAsPending).Methods("POST")

	//Public
	r.HandleFunc("/user", c.CreateAccount).Methods("POST")
	r.HandleFunc("/user/login", c.Authenticate).Methods("POST")

	r.HandleFunc("/address", c.AddAddress).Methods("POST")
	r.HandleFunc("/address", c.DeleteAddress).Methods("DELETE")
	r.HandleFunc("/addresses", c.GetAddresses).Methods("GET")

	r.HandleFunc("/webhooks", c.AddWebhook).Methods("POST")
	r.HandleFunc("/webhooks", c.DeleteWebhook).Methods("DELETE")
	r.HandleFunc("/webhooks", c.GetWebhooks).Methods("GET")

	r.HandleFunc("/invoice", c.CreateInvoice).Methods("POST")
	r.HandleFunc("/invoice", c.UpdateInvoice).Methods("PUT")
	r.HandleFunc("/invoice", c.DeleteInvoice).Methods("DELETE")
	r.HandleFunc("/invoice", c.GetInvoice).Methods("GET")
	r.HandleFunc("/invoices", c.GetInvoices).Methods("GET")

	r.HandleFunc("/health", c.Health).Methods("GET")

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
