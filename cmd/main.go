package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	db = initDB()

	http.Handle("/accounts", http.HandlerFunc(SaveAccount))
	log.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		"localhost", "root", "himanshu@123", "bank", "5432", "disable", "Asia/Kolkata")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("error initializing db %v", err.Error()))
	}
	log.Printf("db started")

	return db
}

type ApiResponse struct {
	Success string
	Message string
}

type AccountReq struct {
	AccountID      int64  `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}

type Account struct {
	ID      int64 `gorm:"column:id;primaryKey"`
	Balance float64
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

	var acc AccountReq
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
	err = SaveAcc(acc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&ApiResponse{
			Success: "false",
			Message: fmt.Sprintf("internal server error: ", err.Error()),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ApiResponse{
		Success: "true",
		Message: "account saved successfully",
	})
}

func SaveAcc(accReq AccountReq) error {
	bal, err := strconv.ParseFloat(accReq.InitialBalance, 64)
	if err != nil {
		return err
	}

	acc := &Account{
		ID:      accReq.AccountID,
		Balance: bal,
	}

	dbResp := db.Create(acc)
	return dbResp.Error
}
