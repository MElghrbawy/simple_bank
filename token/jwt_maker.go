package token

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker returns a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new  token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        payload.ID,
		"username":  payload.Username,
		"issuedAt":  payload.IssuedAt.Unix(),
		"ExpiresAt": payload.ExpiresAt.Unix(),
	})

	tokenString, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, err
	}

	return tokenString, payload, nil
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secretKey := []byte(maker.secretKey)
		return secretKey, nil
	}

	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, ErrInvalidToken
	}

	id, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		return nil, fmt.Errorf("invalid ID in token: %w", err)
	}

	return &Payload{
		ID:        id,
		Username:  claims["username"].(string),
		IssuedAt:  time.Unix(int64(claims["issuedAt"].(float64)), 0),
		ExpiresAt: time.Unix(int64(claims["ExpiresAt"].(float64)), 0),
	}, nil

}
