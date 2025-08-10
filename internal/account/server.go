package account

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Server struct {
	account Accounter
}

type Opt func(*Server)

func WithAccountService(accSvc Accounter) Opt {
	return func(s *Server) {
		s.account = accSvc
	}
}

func NewServer(opt ...Opt) *Server {
	s := &Server{}

	for _, opt := range opt {
		opt(s)
	}

	return s
}

type CreateReq struct {
	ID             uint   `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}

type CreateResp struct {
	Success string `json:"success"`
}

func (r CreateReq) transfomToDTO() (*CreateReqDTO, error) {
	balance, err := decimal.NewFromString(r.InitialBalance)
	if err != nil {
		return nil, err
	}
	return &CreateReqDTO{
		ID:      r.ID,
		Balance: balance,
	}, nil
}

type GetAccountResp struct {
	ID      uint   `json:"id"`
	Balance string `json:"balance"`
}

func (g *GetAccountResp) PopulateFrom(resp *GetAccountDTO) {
	g.ID = resp.ID
	g.Balance = resp.Balance.String()
}

type TransactionReq struct {
	SourceAccountID      uint   `json:"source_account_id"`
	DestinationAccountID uint   `json:"destination_account_id"`
	Amount               string `json:"amount"`
}

func (r *TransactionReq) TransformToDTO() (*TransactionReqDTO, error) {
	amount, err := decimal.NewFromString(r.Amount)
	if err != nil {
		return nil, err
	}

	return &TransactionReqDTO{
		SourceAccountID:      r.SourceAccountID,
		DestinationAccountID: r.DestinationAccountID,
		Amount:               amount,
	}, nil
}

type TransactionResp struct {
	Success string
}

func (s *Server) CreateAccount(c *gin.Context) {
	var req CreateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		log.Printf("failed to bind json: %s", err.Error())
		return
	}

	log.Printf("create account request received for ID: %d", req.ID)

	dto, err := req.transfomToDTO()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("account creation failed: %s", err.Error())
		return
	}

	err = s.account.Create(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		log.Printf("account creation failed: %s", err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (s *Server) GetAccount(c *gin.Context) {
	accID := c.Param("account_id")
	accID = strings.TrimSpace(accID)

	log.Printf("get account request received for ID: %s", accID)

	dto, err := s.account.Get(c.Request.Context(), accID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { 
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account"})
		log.Printf("failed to fetch account: %s", err.Error())
		return
	}

	var resp GetAccountResp
	resp.PopulateFrom(dto)
	c.JSON(http.StatusOK, resp)
}

func (s *Server) Transact(c *gin.Context) {
	var req TransactionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		log.Printf("failed to bind json: %s", err.Error())
		return
	}

	log.Printf("transaction request received: %+v", req)

	dto, err := req.TransformToDTO()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("failed to transform to DTO: %s", err.Error())
		return
	}

	if err := s.account.Transact(c.Request.Context(), dto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		log.Printf("transaction failed: %s", err.Error())
		return
	}

	c.Status(http.StatusOK)
}
