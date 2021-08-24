package model

import "github.com/golang-jwt/jwt"

type Claim struct {
	User `json:"user"`
	jwt.StandardClaims
}
