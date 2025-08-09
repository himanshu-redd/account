package account

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
	ID             int64  `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
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
	ID      int64  `json:"id"`
	Balance string `json:"balance"`
}

func (g *GetAccountResp) PopulateFrom(resp *GetAccountDTO) {
	g.ID = resp.ID
	g.Balance = resp.Balance
}

func (s *Server) CreateAccount(c *gin.Context) {
	var req CreateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("failed to bin json. %s", err.Error())
		c.Error(err)
		return
	}

	log.Printf("create request received")

	dto := req.transfomToDTO()

	err := s.account.Create(c.Request.Context(), dto)
	if err != nil {
		log.Printf("account creation failed. %s", err.Error())
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, CreateResp{
		Success: "true",
	})
}

func (s *Server) GetAccount(c *gin.Context) {
	accID := c.Param("account_id")
	accID = strings.TrimSpace(accID)

	if len(accID) == 0 {
		err := fmt.Errorf("account id should not be empty")
		log.Print(err.Error())
		c.Error(err)
		return
	}

	dto, err := s.account.Get(c.Request.Context(), accID)
	if err != nil {
		log.Printf("failed to fetch resp: %d", err.Error())
		c.Error(err)
		return
	}

	var resp GetAccountResp
	resp.PopulateFrom(dto)
	c.JSON(http.StatusOK, resp)
}
