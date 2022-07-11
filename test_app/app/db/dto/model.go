package dbmodels

import "github.com/golang-jwt/jwt"

// UserDetail Table structure
type UserAuth struct {
	Aud    string `json:"aud"`
	Exp    int    `json:"exp"`
	Sub    string `json:"sub"`
	Email  string `json:"email"`
	UserId int    `json:"id"`
}

type Token struct {
	Email  string `json:"email"`
	UserId string    `json:"id"`
	*jwt.StandardClaims
}
