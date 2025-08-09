package account

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Accounter interface {
	Create(context.Context, *CreateReqDTO) error
	Get(context.Context, string) (*GetAccountDTO, error)
	Transact(context.Context, *TransactionReqDTO) error
}

type CreateReqDTO struct {
	ID      int64
	Balance decimal.Decimal
}

type GetAccountDTO struct {
	ID      int64
	Balance decimal.Decimal
}

type AccountService struct {
	AccountRepo Repoer
}

func NewAccountService(repo Repoer) *AccountService {
	return &AccountService{
		AccountRepo: repo,
	}
}

func (s *AccountService) Create(ctx context.Context, req *CreateReqDTO) error {
	log.Printf("req: %+v", req)

	err := req.validate()
	if err != nil {
		return err
	}

	acc, err := s.AccountRepo.Get(ctx, req.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	log.Printf("creating account to db")

	if acc != nil {
		return fmt.Errorf("account already present")
	}

	if err := s.AccountRepo.Save(ctx, Account{
		ID:      req.ID,
		Balance: req.Balance,
	}); err != nil {
		return err
	}

	return nil
}

func (r *CreateReqDTO) validate() error {
	var errs error

	if r.ID == 0 {
		errs = errors.Join(errs, fmt.Errorf("request id should not be empty"))
	}

	if r.Balance.LessThan(decimal.Zero) {
		errs = errors.Join(errs, fmt.Errorf("invalid initial balance"))
	}

	return errs
}

func (s *AccountService) Get(ctx context.Context, accountID string) (*GetAccountDTO, error) {
	log.Printf("fetch account details request received for account id: %s", accountID)

	if len(accountID) == 0 {
		err := fmt.Errorf("account id should not be empty")
		log.Print(err.Error())
		return nil, err
	}

	accID, err := strconv.ParseInt(accountID, 10, 64)
	if err != nil {
		return nil, err
	}

	accDetails, err := s.AccountRepo.Get(ctx, accID)
	if err != nil {
		return nil, err
	}

	return &GetAccountDTO{
		ID:      accDetails.ID,
		Balance: accDetails.Balance,
	}, nil
}

type TransactionReqDTO struct {
	SourceAccountID      int64
	DestinationAccountID int64
	Amount               decimal.Decimal
}

func (r *TransactionReqDTO) validate() error {
	var errs error

	if r.SourceAccountID == 0 {
		errs = errors.Join(errs, fmt.Errorf("source account should not be empty"))
	}
	if r.DestinationAccountID == 0 {
		errs = errors.Join(errs, fmt.Errorf("destination account should not be empty"))
	}
	if r.Amount.LessThan(decimal.Zero) {
		errs = errors.Join(errs, fmt.Errorf("deduction amount should be positive"))
	}

	return errs
}

func (s *AccountService) Transact(ctx context.Context, req *TransactionReqDTO) error {
	if err := req.validate(); err != nil {
		return err
	}

	sourceAcc, err := s.AccountRepo.Get(ctx, req.SourceAccountID)
	if err != nil {
		return err
	}

	destinationAcc, err := s.AccountRepo.Get(ctx, req.DestinationAccountID)
	if err != nil {
		return err
	}

	if sourceAcc.Balance.LessThan(req.Amount) {
		err := fmt.Errorf("insufficient balance")
		fmt.Printf("%s", err.Error())
		return err
	}

	sourceAcc.Balance = sourceAcc.Balance.Sub(req.Amount)
	destinationAcc.Balance = destinationAcc.Balance.Add(req.Amount)

	s.AccountRepo.Save(ctx, *sourceAcc)
	s.AccountRepo.Save(ctx, *destinationAcc)

	return nil
}
