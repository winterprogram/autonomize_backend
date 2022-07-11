package jwt

import (
	"context"
	"fmt"
	dbmodels "test/test_app/app/db/dto"
	"test/test_app/app/service/util"

	"github.com/dgrijalva/jwt-go"
)

type IJwtService interface {
	VerifyToken(ctx context.Context, tokenString string) (*dbmodels.Token, bool)
}

type JwtService struct{}

func NewJwtService() IJwtService {
	return &JwtService{}
}
func (j *JwtService) VerifyToken(ctx context.Context, tokenString string) (*dbmodels.Token, bool) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(util.GetEnvWithKey("JWT_ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, false
	}

	var auth dbmodels.Token
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		auth.UserId = claims["id"].(string)
		auth.Email = claims["email"].(string)
		return &auth, true
	}
	return nil, false

}
