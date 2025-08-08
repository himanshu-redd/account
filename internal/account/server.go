package account

import (
	"log"
	"net/http"

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

func (s *Server) CreateAccount(c *gin.Context) {
	var req CreateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("error: failed to bin json. %s", err.Error())
		c.Error(err)
		return
	}

	log.Printf("create request received")

	dto := req.transfomToDTO()

	err := s.account.Create(c.Request.Context(), dto)
	if err != nil {
		log.Printf("error: account creation failed. %s", err.Error())
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, CreateResp{
		Success: "true",
	})
}
