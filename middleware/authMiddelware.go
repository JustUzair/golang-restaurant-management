package middleware

import (
	"net/http"
	"restaurant-management/utils"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt, _ := c.Cookie("jwt")
		// jwtRefresh, _ := c.Cookie("jwt_refresh")
		clientToken := jwt
		if jwt == "" {
			clientToken = c.Request.Header.Get("token")
		}
		if clientToken == "" {
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    "table updation failed",
			})
			return
		}
		claims, err := utils.ValidateToken(clientToken)
		if err != "" {
			errorCode := http.StatusInternalServerError
			c.JSON(errorCode, gin.H{
				"statusCode": errorCode,
				"status":     "error",
				"message":    err,
			})
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.UID)
		c.Next()
	}
}
