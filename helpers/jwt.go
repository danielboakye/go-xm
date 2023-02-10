package helpers

import (
	"time"

	"github.com/danielboakye/go-xm/config"
	"github.com/golang-jwt/jwt/v4"
)

const (
	accessKeyType = "access"
)

type JWTClaims struct {
	UserID  string
	KeyType string
	jwt.RegisteredClaims
}

func GenerateAccessToken(
	cfg config.Configurations,
	userID string,
) (token string, err error) {

	claims := &JWTClaims{
		UserID:  userID,
		KeyType: accessKeyType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Local().Add(cfg.AccessTokenDuration),
			),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JWTSecretKey))
	return
}
func ValidateAccessToken(signedToken string, cfg config.Configurations) (claims *JWTClaims, err error) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return []byte(cfg.JWTSecretKey), nil
		},
	)

	if err != nil {
		err = ErrInvalidToken
		return
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != accessKeyType {
		err = ErrUnauthorized
		return
	}

	return
}
