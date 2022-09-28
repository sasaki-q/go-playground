package factory

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeyLength = 12

var (
	invalidErr = errors.New("ERROR: token is invalid")
	expiredErr = errors.New("ERROR: token is expired")
)

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return expiredErr
	}

	return nil
}

type JwtTokenFactory struct {
	secretKey string
}

func GenerateTokenFactory(secretKey string) (Factory, error) {
	if len(secretKey) < minSecretKeyLength {
		return nil, fmt.Errorf("ERROR: invalid key length, min length at %d", minSecretKeyLength)
	}

	return &JwtTokenFactory{secretKey}, nil
}

func (factory *JwtTokenFactory) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := GeneratePayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(factory.secretKey))
}

func (factory *JwtTokenFactory) VerifyToken(token string) (*Payload, error) {
	keyFn := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, invalidErr
		}

		return []byte(factory.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFn)

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, expiredErr) {
			return nil, expiredErr
		}

		return nil, invalidErr
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, invalidErr
	}

	return payload, nil
}
