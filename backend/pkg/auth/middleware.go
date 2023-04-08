package auth

import (
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Middleware checks if the JWT sent is valid or not. This function is involked for every API route that needs authentication
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, "Authorization Token is required")
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}

		extractedToken := strings.Split(clientToken, "Bearer ")

		// Verify if the format of the token is correct
		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			c.JSON(http.StatusBadRequest, "Incorrect Format of Authorization Token")
			c.Abort()
			return
		}

		// Parse the claims
		parsedToken, err := jwt.ParseWithClaims(clientToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("ACCESS_SECRET")), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, "Invalid Token Signature")
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, "Bad Request")
			c.Abort()
			return
		}

		if !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, "Invalid Token")
			c.Abort()
			return
		}

		// currently only admin has permission
		role := claims["role"].(string)
		/*if role != "admin" {
			c.JSON(http.StatusUnauthorized, "Permission denied")
			c.Abort()
			return
		}*/

		c.Set("id", claims["id"].(string))
		c.Set("role", role)

		c.Next()
	}
}
