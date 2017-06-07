package server

import (
	"fmt"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// Authenticate auth by HTTP Method POST
func (env *Env) Authenticate(c *gin.Context) {
	var requestJSON struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	log.Println(c.Request)
	if c.BindJSON(&requestJSON) == nil {
		log.Println("requestJSON :", requestJSON)
		user, err := env.db.GetUserAccount(requestJSON.Account)
		if err != nil {
			errMsg := fmt.Sprint("Get User failed with", err)
			c.JSON(500, gin.H{"message": errMsg})
			return
		}
		if user.Password == requestJSON.Password {
			log.Println("user :", user)
			// Create the token
			token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
			// Set some claims
			token.Claims = jwt_lib.MapClaims{
				"Account": user.Account,
				"Role":    user.Role,
				"exp":     time.Now().Add(time.Hour * 12).Unix(),
			}
			// Sign and get the complete encoded token as a string
			tokenString, err := token.SignedString([]byte(env.jwtSecret))
			if err != nil {
				c.JSON(500, gin.H{"message": "Could not generate token"})
			}
			c.JSON(200, gin.H{"token": tokenString})
			return
		}
		c.JSON(401, gin.H{"message": "Wrong account or password"})
		return
	}
	c.JSON(401, gin.H{"message": "Invalid Request"})
}