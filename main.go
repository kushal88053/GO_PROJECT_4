package main

import (
	"encoding/json"
	"net"
	"net/http"
	"regexp"
	"strings"
)

type EmailResponse struct {
	Email     string `json:"email"`
	IsValid   bool   `json:"is_valid_format"`
	DomainOK  bool   `json:"domain_has_mx"`
	FullValid bool   `json:"fully_valid"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func validateEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	w.Header().Set("Content-Type", "application/json")

	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "email parameter is required"})
		return
	}

	isValid := emailRegex.MatchString(email)
	domainOK := false

	if isValid {
		domain := strings.Split(email, "@")[1]
		mxRecords, err := net.LookupMX(domain)
		if err == nil && len(mxRecords) > 0 {
			domainOK = true
		}
	}

	resp := EmailResponse{
		Email:     email,
		IsValid:   isValid,
		DomainOK:  domainOK,
		FullValid: isValid && domainOK,
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/validate", validateEmailHandler)
	println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
