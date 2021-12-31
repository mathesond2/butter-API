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
	r.HandleFunc("/api/webhooks/mempoolEvent", c.ParseMempoolEvent).Methods("POST")
	r.HandleFunc("/api/webhooks/updateInvoiceStatusFromEvent", c.UpdateInvoiceStatusFromEvent).Methods("POST") //prob should be put
	r.HandleFunc("/api/invoice/status", c.UpdateInvoiceStatusAsPending).Methods("POST")

	//Public
	r.HandleFunc("/api/user", c.CreateAccount).Methods("POST")
	r.HandleFunc("/api/user/login", c.Authenticate).Methods("POST")

	r.HandleFunc("/api/address", c.AddAddress).Methods("POST")
	r.HandleFunc("/api/address", c.DeleteAddress).Methods("DELETE")
	r.HandleFunc("/api/addresses", c.GetAddresses).Methods("GET")

	r.HandleFunc("/api/webhooks", c.AddWebhook).Methods("POST")
	r.HandleFunc("/api/webhooks", c.DeleteWebhook).Methods("DELETE")
	r.HandleFunc("/api/webhooks", c.GetWebhooks).Methods("GET")

	r.HandleFunc("/api/invoice", c.CreateInvoice).Methods("POST")
	r.HandleFunc("/api/invoice", c.UpdateInvoice).Methods("PUT")
	r.HandleFunc("/api/invoice", c.DeleteInvoice).Methods("DELETE")
	r.HandleFunc("/api/invoice", c.GetInvoice).Methods("GET")
	r.HandleFunc("/api/invoices", c.GetInvoices).Methods("GET")

	r.HandleFunc("/api/health", c.Health).Methods("GET")

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
