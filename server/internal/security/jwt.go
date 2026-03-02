package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret    []byte
	issuer    string
	accessTTL time.Duration
}

func NewTokenManager(secret, issuer string, accessTTL time.Duration) (*TokenManager, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is empty")
	}
	if issuer == "" {
		issuer = "server"
	}
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}

	return &TokenManager{
		secret:    []byte(secret),
		issuer:    issuer,
		accessTTL: accessTTL,
	}, nil
}

func (tm *TokenManager) GenerateAccessToken(userID string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tm.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tm.accessTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.secret)
}

func (tm *TokenManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected jwt signing method")
		}
		return tm.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
