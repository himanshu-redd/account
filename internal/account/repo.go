package account

import (
	"context"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type Repoer interface {
	Create(context.Context, *CreateReqDTO) error
}

type AccountRepo struct {
	DB *gorm.DB
}

func NewAccountRepo(db *gorm.DB) *AccountRepo {
	return &AccountRepo{
		DB: db,
	}
}

type Account struct {
	ID      int64   `gorm:"column:id;primaryKey"`
	Balance float64 `gorm:"column:balance"`
}

func (a *Account) populateFrom(req *CreateReqDTO) error {
	a.ID = req.ID
	var err error
	a.Balance, err = strconv.ParseFloat(req.InitialBalance, 64)
	if err != nil {
		return err
	}
	return nil
}

func (r *AccountRepo) Create(ctx context.Context, req *CreateReqDTO) error {
	var acc Account

	if err := acc.populateFrom(req); err != nil {
		return err
	}
	log.Printf("account model: %+v", acc)

	dbResp := r.DB.Create(&acc)
	if dbResp.Error != nil {
		log.Printf("error occurred while creating account: %s", dbResp.Error.Error())
		return dbResp.Error
	}

	return nil
}
