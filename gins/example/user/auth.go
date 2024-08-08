package user

import (
	"fmt"
	"github.com/aiechoic/services/openapi"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var BearerSecurityScheme = "bearerAuth"

var SecuritySchemes = openapi.SecuritySchemes{
	BearerSecurityScheme: &openapi.SecurityScheme{
		Type: "apiKey",
		In:   "header",
		Name: "Authorization",
	},
	// Add more security schemes here
}

type JWTAuth struct {
	secretKey string
}

func NewJWTAuth(secretKey string) *JWTAuth {
	return &JWTAuth{
		secretKey: secretKey,
	}
}

func (j *JWTAuth) GetToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func (j *JWTAuth) Auth(c *gin.Context) {
	token := j.GetToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		c.Abort()
		return
	}

	user, err := j.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}

func (j *JWTAuth) GenerateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"id":   user.Id,
		"name": user.Name,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTAuth) ParseToken(tokenString string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &User{
			Id:   int(claims["id"].(float64)),
			Name: claims["name"].(string),
		}, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func (j *JWTAuth) GetUser(c *gin.Context) *User {
	user, ok := c.Get("user")
	if !ok {
		token := j.GetToken(c)
		u, err := j.ParseToken(token)
		if err != nil {
			return nil
		}
		c.Set("user", u)
		return u
	}
	return user.(*User)
}

func (j *JWTAuth) SecurityScheme() []map[string][]string {
	return []map[string][]string{
		{
			BearerSecurityScheme: {},
		},
	}
}
