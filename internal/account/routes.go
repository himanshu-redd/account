package account

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes(s *Server) {
	router := gin.Default()

	router.POST("/accounts", s.CreateAccount)
	router.GET("/accounts/:account_id", s.GetAccount)
	router.POST("/transactions", s.Transact)

	router.Run(":8080")
}
