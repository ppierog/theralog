package cred

import (
	"github.com/golang-jwt/jwt"
)

// https://pascalallen.medium.com/jwt-authentication-with-go-242215a9b4f8
type UserClaims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

type Token struct {
	Jwt string `json:"jwt"`
}

type TokenRepository struct {
	Secret string
}

func (p *TokenRepository) NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(p.Secret))
}

func (p *TokenRepository) NewRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(p.Secret))
}

func (p *TokenRepository) ParseAccessToken(accessToken string) *UserClaims {
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.Secret), nil
	})

	return parsedAccessToken.Claims.(*UserClaims)
}

func (p *TokenRepository) ParseRefreshToken(refreshToken string) *jwt.StandardClaims {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.Secret), nil
	})

	return parsedRefreshToken.Claims.(*jwt.StandardClaims)
}
