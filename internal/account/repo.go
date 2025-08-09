package account

import (
	"context"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type Repoer interface {
	Save(context.Context, *CreateReqDTO) error
	Get(context.Context, int64) (*GetAccountDTO, error)
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
	gorm.Model
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

func (a *Account) Transform() *GetAccountDTO {
	return &GetAccountDTO{
		ID:      a.ID,
		Balance: a.Balance,
	}
}

func (r *AccountRepo) Save(ctx context.Context, req *CreateReqDTO) error {
	var acc Account

	if err := acc.populateFrom(req); err != nil {
		return err
	}
	log.Printf("account model: %+v", acc)

	dbResp := r.DB.Save(&acc)
	if dbResp.Error != nil {
		log.Printf("error occurred while creating account: %s", dbResp.Error.Error())
		return dbResp.Error
	}

	return nil
}

func (r *AccountRepo) Get(ctx context.Context, accID int64) (*GetAccountDTO, error) {
	var acc Account

	dbResp := r.DB.First(&acc, accID)
	if dbResp.Error != nil {
		log.Printf("error occurred while fetching account:  %s", dbResp.Error.Error())
		return nil, dbResp.Error
	}

	dto := acc.Transform()
	return dto, nil
}
