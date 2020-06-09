package middleware

import (
	"fmt"
	"hangout-api/context"
	"hangout-api/header"
	"hangout-api/models"

	"github.com/gin-gonic/gin"
)

func RequireUser(us models.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := header.GetToken(c)
		if token == "" {
			c.JSON(401, gin.H{
				"message": "token null",
			})
			fmt.Println("token null")
			c.Abort()
			return
		}
		user, err := us.GetByToken(token)
		if err != nil {
			c.JSON(401, gin.H{
				"message": "can't get by token",
			})
			fmt.Println("can't get by token")
			c.Abort()
			return
		}
		context.SetUser(c, user)
	}
}
