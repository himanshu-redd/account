package account

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Accounter interface {
	Create(context.Context, *CreateReqDTO) error
	Get(context.Context, string) (*GetAccountDTO, error)
}

type CreateReqDTO struct {
	ID             int64
	InitialBalance string
}

type GetAccountDTO struct {
	ID      int64
	Balance string
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

	err := req.validateCreateReq()
	if err != nil {
		return err
	}

	if err := s.AccountRepo.Create(ctx, req); err != nil {
		return err
	}

	return nil
}

func (r *CreateReqDTO) validateCreateReq() error {
	var errs error

	if r.ID == 0 {
		errs = errors.Join(errs, fmt.Errorf("request id should not be empty"))
	}
	if len(strings.TrimSpace(r.InitialBalance)) == 0 {
		errs = errors.Join(errs, fmt.Errorf("initla balance should not be empty"))
	} else {

		balance, err := strconv.ParseFloat(r.InitialBalance, 64)

		if err != nil {
			errs = errors.Join(errs, err)
		} else if balance < 0 {
			errs = errors.Join(err, fmt.Errorf("invalid initial balance"))
		}
	}

	return errs
}

func (s *AccountService) Get(ctx context.Context, accountID string) (*GetAccountDTO, error) {
	log.Printf("Request received for account id: %s", accountID)

	accID, err := strconv.ParseInt(accountID, 10, 64)
	if err != nil {
		return nil, err
	}

	accDetails, err := s.AccountRepo.Get(ctx, accID)
	if err != nil {
		return nil, err
	}

	return accDetails, err
}
