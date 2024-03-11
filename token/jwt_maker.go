package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

// JwtMaker is a JSON Web Token Maker
type JwtMaker struct {
	secretKey string
}

// NewJWTMaker create a new jwtmaker
func CreateJWTMaker(secret string) (Maker, error) {
	if len(secret) < minSecretKeySize {
		return nil, ErrInvalidLenSecretKey
	}
	return &JwtMaker{secret}, nil
}

// Createtoken create new token from username and duration
func (maker *JwtMaker) CreateToken(username string, duration time.Duration) (string, error) {
	fmt.Println("newPayload")
	payload, err := NewPayload(username, duration)
	fmt.Println(payload)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	fmt.Println(jwtToken)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken check if token valid or not
func (maker *JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpireToken) {
			return nil, ErrExpireToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil

}
