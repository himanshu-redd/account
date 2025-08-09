package account

import (
	"context"
	"log"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Repoer interface {
	Save(context.Context, Account) error
	Get(context.Context, int64) (*Account, error)
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
	ID      int64           `gorm:"column:id;primaryKey"`
	Balance decimal.Decimal `gorm:"column:balance"`
}

func (a *Account) Transform() *GetAccountDTO {
	return &GetAccountDTO{
		ID:      a.ID,
		Balance: a.Balance,
	}
}

func (r *AccountRepo) Save(ctx context.Context, account Account) error {
	dbResp := r.DB.Save(&account)
	if dbResp.Error != nil {
		log.Printf("error occurred while creating account: %s", dbResp.Error.Error())
		return dbResp.Error
	}

	return nil
}

func (r *AccountRepo) Get(ctx context.Context, accID int64) (*Account, error) {
	var acc Account

	dbResp := r.DB.First(&acc, accID)
	if dbResp.Error != nil {
		log.Printf("error occurred while fetching account:  %s", dbResp.Error.Error())
		return nil, dbResp.Error
	}

	return &acc, nil
}
