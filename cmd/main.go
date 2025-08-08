package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.Handle("/accounts", http.HandlerFunc(SaveAccount))
	log.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


type ApiResponse struct {
	Success string
	Message string
}

type Account struct {
	AccountID      int64  `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}

func SaveAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ApiResponse{
			Success: "false",
			Message: "only post method allowd",
		})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ApiResponse{
			Success: "false",
			Message: fmt.Sprintf("invalid request body, err: ", err.Error()),
		})
		return
	}

	var acc Account
	err = json.Unmarshal(body, &acc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ApiResponse{
			Success: "false",
			Message: fmt.Sprintf("invalid JSON body. err: %s", err.Error()),
		})
		return
	}

	fmt.Printf("body: %+v", acc)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ApiResponse{
		Success: "true",
		Message: "account saved successfully",
	})
}
