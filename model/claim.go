package model

import "github.com/golang-jwt/jwt"

type Claim struct {
	User `json:"user"`
	jwt.StandardClaims
}

type ClaimRefreshToken struct {
	Username    string `json:"username"`
	RefreshUuid string `json:"RefreshUuid"`
	jwt.StandardClaims
}
