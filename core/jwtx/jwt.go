package jwtx

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken generate tokens used for auth
func GenerateToken(secret string, claims jwt.Claims) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func DecodeToken(secret string, token string, skipValidation bool, claim jwt.Claims) error {
	p := new(jwt.Parser)
	p.SkipClaimsValidation = false
	tkn, err := p.ParseWithClaims(token, claim, func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return fmt.Errorf("invalid token, err: %v", err)
	}
	if !tkn.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}
