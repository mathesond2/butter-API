package main

import (
	"fmt"
	"go-invoices/app"
	c "go-invoices/controllers"
	u "go-invoices/utils"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
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

	credentials := handlers.AllowCredentials()
	//no AllowedHeaders for now?
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS"})
	// ttl := handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins([]string{"www.justbutter.co", "http://localhost:3000", "http://localhost:3000/pay/3", "http://justbutter.co", "http://www.justbutter.co"})

	err := http.ListenAndServe(":"+port, handlers.CORS(credentials, allowedHeaders, methods, origins)(r))
	if err != nil {
		fmt.Print(err)
	}
}
