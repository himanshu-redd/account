package account

import (
	"context"
	"fmt"
	"log"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repoer interface {
	Create(context.Context, Account) error
	Get(context.Context, uint) (*Account, error)
	Transfer(context.Context, uint, uint, decimal.Decimal) error
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
	Balance decimal.Decimal `gorm:"column:balance"`
}

func (a *Account) Transform() *GetAccountDTO {
	return &GetAccountDTO{
		ID:      a.ID,
		Balance: a.Balance,
	}
}

func (r *AccountRepo) Create(ctx context.Context, account Account) error {
	dbResp := r.DB.WithContext(ctx).Create(&account)
	if dbResp.Error != nil {
		log.Printf("error occurred while creating account: %s", dbResp.Error.Error())
		return dbResp.Error
	}
	return nil
}

func (r *AccountRepo) Get(ctx context.Context, accID uint) (*Account, error) {
	var acc Account

	dbResp := r.DB.First(&acc, accID)
	if dbResp.Error != nil {
		log.Printf("error occurred while fetching account:  %s", dbResp.Error.Error())
		return nil, dbResp.Error
	}

	return &acc, nil
}

func (r *AccountRepo) Transfer(ctx context.Context, fromID, toID uint, amount decimal.Decimal) error {
	return r.DB.WithContext(ctx).Transaction(func(db *gorm.DB) error {
		var fromAccount Account
		if err := db.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&fromAccount, fromID).Error; err != nil {
			return err
		}

		if fromAccount.Balance.LessThan(amount) {
			return fmt.Errorf("insufficient balance in account ID %d", fromID)
		}

		var toAccount Account
		if err := db.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&toAccount, toID).Error; err != nil {
			return err
		}

		fromAccount.Balance = fromAccount.Balance.Sub(amount)
		toAccount.Balance = toAccount.Balance.Add(amount)

		if err := db.Save(&fromAccount).Error; err != nil {
			return err
		}
		if err := db.Save(&toAccount).Error; err != nil {
			return err
		}

		return nil
	})
}
