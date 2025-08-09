package account

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
	ID             int64   `json:"account_id"`
	InitialBalance float64 `json:"initial_balance"`
}

type CreateResp struct {
	Success string `json:"success"`
}

func (r CreateReq) transfomToDTO() *CreateReqDTO {
	return &CreateReqDTO{
		ID:             r.ID,
		InitialBalance: r.InitialBalance,
	}
}

type GetAccountResp struct {
	ID      int64   `json:"id"`
	Balance float64 `json:"balance"`
}

func (g *GetAccountResp) PopulateFrom(resp *GetAccountDTO) {
	g.ID = resp.ID
	g.Balance = resp.Balance
}

type TransactionReq struct {
	SourceAccountID      int64   `json:"source_account_id"`
	DestinationAccountID int64   `json:"destination_account_id"`
	Amount               float64 `json:"amount"`
}

func (r *TransactionReq) TransformToDTO() *TransactionReqDTO {
	return &TransactionReqDTO{
		SourceAccountID:      r.SourceAccountID,
		DestinationAccountID: r.DestinationAccountID,
		Amount:               r.Amount,
	}
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

	dto := req.transfomToDTO()

	err := s.account.Create(c.Request.Context(), dto)
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

	log.Printf("get account request received for ID: %d", accID)

	dto, err := s.account.Get(c.Request.Context(), accID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // Example of a specific error message
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

	dto := req.TransformToDTO()
	if err := s.account.Transact(c.Request.Context(), dto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		log.Printf("transaction failed: %s", err.Error())
		return
	}

	c.Status(http.StatusOK)
}
