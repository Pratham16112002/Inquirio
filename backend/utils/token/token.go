package token

import (
	"Inquiro/config/env"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

var JWT *JWTAuthenticator

func init() {
	JWT = NewJWT(
		env.GetString("JWT_SECRET", "1234"),
		env.GetString("JWT_AUD", "inquirio"),
		env.GetString("JWT_ISS", "inquirio"),
	)
}

func NewJWT(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret: secret,
		aud:    aud,
		iss:    iss,
	}
}

func (j *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte((j.secret)))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.secret), nil
	}, jwt.WithAudience(j.aud), jwt.WithExpirationRequired(), jwt.WithIssuer(j.iss), jwt.WithValidMethods([]string{jwt.SigningMethodES256.Name}))
}
