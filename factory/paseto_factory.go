package factory

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoFactory struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func GeneratePasetoFactory(key string) (Factory, error) {
	if len(key) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("ERROR: invalid key length, min length at %d", chacha20poly1305.KeySize)
	}

	factory := &PasetoFactory{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(key),
	}

	return factory, nil
}

func (factory *PasetoFactory) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := GeneratePayload(username, duration)
	if err != nil {
		return "", err
	}

	return factory.paseto.Encrypt(factory.symmetricKey, payload, nil)
}

func (factory *PasetoFactory) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := factory.paseto.Decrypt(token, factory.symmetricKey, payload, nil)
	if err != nil {
		return nil, invalidErr
	}

	err = payload.Valid()
	if err != nil {
		return nil, expiredErr
	}

	return payload, nil
}
