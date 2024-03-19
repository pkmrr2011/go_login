package middleware

import (
	"errors"
	"fmt"
	"login/shared"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	fmt.Println("==", token)
	return nil
}

func ExtractBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is missing")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Authorization header is not correctly formatted")
	}
	return parts[1], nil
}

func Authorization(c *gin.Context) {
	token, err := ExtractBearerToken(c)
	if err != nil {
		shared.HandleError(c, http.StatusUnauthorized, "Error extracting token", err)
		return
	}

	err = VerifyToken(token)

	if err != nil {
		shared.HandleError(c, http.StatusUnauthorized, "Error verifying token", err)
		return
	}
	c.Next()
}
