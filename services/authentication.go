package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

const (
	HashCost = 10
	TokenTypeAccess  = "ACCESS_TOKEN"
	TokenTypeRefresh = "REFRESH_TOKEN"
)

var (
	tokenTypeExpiryMap = gin.H{
		TokenTypeAccess:	time.Minute * 60, // 60 minutes
		TokenTypeRefresh:	time.Hour * 24 * 60, // 60 days
	}
)

type HasLoginIdentity interface {
	GetLoginId() string
}

func HashPassword(password string) (string, error) {
	passBytes, err := bcrypt.GenerateFromPassword([]byte(password), HashCost)

	return string(passBytes), err
}

func CompareHashAndPass(hash string, pass string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}

func GenerateNewJwtToken(id HasLoginIdentity, tokenType string) (string, error) {
	iat := jwt.NewNumericDate(time.Now())
	nbf := jwt.NewNumericDate(time.Now())
	exp := jwt.NewNumericDate(time.Now().Add(tokenTypeExpiryMap[tokenType].(time.Duration)))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    id.GetLoginId(),
		Subject:   tokenType,
		ExpiresAt: exp,
		NotBefore: nbf,
		IssuedAt:  iat,
		ID:        uuid.NewString(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_HS256_KEY")))
}

func VerifyJwtToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_HS256_KEY")), nil
	})
	if err != nil {
		log.Println("could not parse JWT Token:", token)
		log.Println(err)
		return nil, err
	}

	if _, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return parsedToken, nil
	}

	return nil, fmt.Errorf("found invalid JWT token: %s", token)
}