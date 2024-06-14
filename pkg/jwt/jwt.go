package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Manager interface {
	NewToken(id string) (string, error)
	ParseToken(tokenString string) (string, error)
}

type ManagerConfig struct {
	SigningKey string
	TokenTtl   time.Duration
}

type manager struct {
	signingKey string
	tokenTtl   time.Duration
}

type CustomClaims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}

func NewManager(config *ManagerConfig) (*manager, error) {
	if config.SigningKey == "" {
		return nil, errors.New("empty signing key")
	}
	return &manager{
		signingKey: config.SigningKey,
		tokenTtl:   config.TokenTtl,
	}, nil
}

func (m *manager) NewToken(id string) (string, error) {
	claims := CustomClaims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.tokenTtl).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(m.signingKey))
	return signedString, err
}

func (m *manager) ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("error getting claims from token")
	}
	return claims["id"].(string), nil
}
