package token

import (
	"testing"
	"time"

	"github.com/MElghrbawy/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	ExpiresAt := issuedAt.Add(duration)

	token, _, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, ExpiresAt, payload.ExpiresAt, time.Second)

}

// func TestExpiredJWTToken(t *testing.T) {
// 	maker, err := NewJWTMaker(util.RandomString(32))
// 	require.NoError(t, err)

// 	token, err := maker.CreateToken(util.RandomOwner(), -time.Second)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, token)

// 	payload, err := maker.VerifyToken(token)
// 	require.Error(t, err)
// 	require.EqualError(t, err, ErrInvalidToken.Error())
// 	require.Nil(t, payload)
// }

// func TestInvalidJWTTokenAlgNone(t *testing.T) {
// 	payload, err := NewPayload(util.RandomOwner(), time.Minute)
// 	require.NoError(t, err)

// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
// 		"id":        payload.ID,
// 		"username":  payload.Username,
// 		"issuedAt":  payload.IssuedAt.Unix(),
// 		"ExpiresAt": payload.ExpiresAt.Unix(),
// 	})

// 	tokenString, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)

// 	require.NoError(t, err)
// 	require.NotEmpty(t, tokenString)

// 	maker, err := NewJWTMaker(util.RandomString(32))
// 	require.NoError(t, err)

// 	payload, err = maker.VerifyToken(tokenString)
// 	require.Error(t, err)
// 	require.EqualError(t, err, ErrInvalidToken.Error())
// 	require.Nil(t, payload)
// }
