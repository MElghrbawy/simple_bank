package token

import (
	"encoding/json"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

func NewPasetoMaker(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size of %d, expected %d", len(symmetricKey), chacha20poly1305.KeySize)
	}

	key := paseto.NewV4SymmetricKey()
	maker := &PasetoMaker{

		symmetricKey: key,
	}

	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	token := paseto.NewToken()

	token.SetIssuedAt(payload.IssuedAt)
	token.SetNotBefore(payload.IssuedAt)
	token.SetExpiration(payload.ExpiresAt)

	token.SetString("id", payload.ID.String())
	token.SetString("username", payload.Username)

	return token.V4Encrypt(maker.symmetricKey, nil), payload, nil

}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(tokenString string) (*Payload, error) {
	payload := &Payload{}

	parser := paseto.NewParser()
	token, err := parser.ParseV4Local(maker.symmetricKey, tokenString, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	jsonClaims := token.ClaimsJSON()

	var claims map[string]interface{}
	if err := json.Unmarshal([]byte(jsonClaims), &claims); err != nil {
		return nil, err // Handle unmarshal error appropriately
	}
	issuedAt, err := time.Parse(time.RFC3339, claims["iat"].(string))
	if err != nil {
		return nil, err
	}
	expires, err := time.Parse(time.RFC3339, claims["exp"].(string))
	if err != nil {
		return nil, err
	}

	id, _ := uuid.Parse(claims["id"].(string))

	payload.ID = id
	payload.Username = claims["username"].(string)
	payload.IssuedAt = issuedAt
	payload.ExpiresAt = expires

	return payload, nil
}
