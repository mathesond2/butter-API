package utils

import (
	"encoding/json"
	"net/http"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func IsValidEthAddress(address string) bool {
	if len(address) != 42 {
		return false
	}

	if address[0:2] != "0x" {
		return false
	}

	return true
}
