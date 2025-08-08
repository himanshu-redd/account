package main

import (
	"github.com/himanshu-redd/valuelabs-assignment/internal/account"
	config "github.com/himanshu-redd/valuelabs-assignment/pkg"
)

func main() {
	db := config.InitDB()
	AccRepo := account.NewAccountRepo(db)
	AccSvc := account.NewAccountService(AccRepo)
	server := account.NewServer(
		account.WithAccountService(AccSvc),
	)
	account.InitRoutes(server)
}
