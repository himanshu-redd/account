package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes(s *Server) {
	router := gin.Default()

	router.Use(ErrorMiddleware())

	router.POST("/accounts", s.CreateAccount)
	router.GET("/accounts/:account_id", s.GetAccount)

	router.Run(":8080")
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err != nil {
			c.JSON(http.StatusBadRequest, &struct {
				Success string `json:"success"`
				Message string `json:"message"`
			}{
				Success: "false",
				Message: err.Err.Error(),
			})
		}
	}
}
