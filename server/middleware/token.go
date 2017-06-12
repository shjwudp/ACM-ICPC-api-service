package middleware

import (
	"errors"
	jwt_lib "github.com/dgrijalva/jwt-go"
	jwt_request "github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"github.com/shjwudp/ACM-ICPC-api-service/model"
)

// FakeJWTAuthMiddleware : for performance test
func FakeJWTAuthMiddleware(c *gin.Context) {
	c.Set("user", model.User{Role: "admin"})
	c.Next()
}

// JWTAuthMiddleware : JWT Authorization Verification
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt_request.ParseFromRequest(
			c.Request,
			jwt_request.OAuth2Extractor,
			func(token *jwt_lib.Token) (interface{}, error) {
				b := ([]byte(secret))
				return b, nil
			},
		)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		c.Set("token", token)
		account, ok := token.Claims.(jwt_lib.MapClaims)["Account"]
		if !ok {
			c.AbortWithError(500, errors.New("No Account in token"))
			return
		}
		role, ok := token.Claims.(jwt_lib.MapClaims)["Role"]
		if !ok {
			c.AbortWithError(500, errors.New("No Role in token"))
			return
		}

		user := model.User{
			Account: account.(string),
			Role:    role.(string),
		}
		c.Set("user", user)
		c.Next()
	}
}

var level2Role = map[string]bool{
	"normal": true,
}

var level1Role = map[string]bool{
	"volunteer": true,
}

var level0Role = map[string]bool{
	"admin": true,
}

// Level1PermissionMiddleware only access <= level1 role
func Level1PermissionMiddleware(c *gin.Context) {
	raw, has := c.Get("user")
	if !has {
		c.AbortWithError(500, errors.New("No user in the gin.Context"))
		return
	}
	user := raw.(model.User)
	if _, ok := level0Role[user.Role]; ok {
		c.Next()
		return
	}
	if _, ok := level1Role[user.Role]; ok {
		c.Next()
		return
	}
	c.AbortWithError(403, errors.New("No Permission"))
}

// Level0PermissionMiddleware only access <= level0 top
func Level0PermissionMiddleware(c *gin.Context) {
	raw, has := c.Get("user")
	if !has {
		c.AbortWithError(500, errors.New("No user in the gin.Context"))
		return
	}
	user := raw.(model.User)
	if _, ok := level0Role[user.Role]; ok {
		c.Next()
		return
	}
	c.AbortWithError(403, errors.New("No Permission"))
}
