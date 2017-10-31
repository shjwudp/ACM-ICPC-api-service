package server

import (
	"database/sql"
	"fmt"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// Authenticate auth by HTTP Method POST
func (env *Env) Authenticate(c *gin.Context) {
	var requestJSON struct {
		Account  string `json:"account" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	log.Println(c.Request)
	err := c.BindJSON(&requestJSON)
	if err != nil {
		errMsg := fmt.Sprint("BindJSON failed with", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": errMsg})
		return
	}
	log.Println("requestJSON :", requestJSON)
	user, err := env.db.GetUserAccount(requestJSON.Account)
	if err == sql.ErrNoRows {
		errMsg := fmt.Sprint("Get User failed with", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": errMsg})
		return
	} else if err != nil {
		errMsg := fmt.Sprint("Get User failed with", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
		return
	}
	if user.Password != requestJSON.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong account or password"})
		return
	}
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
