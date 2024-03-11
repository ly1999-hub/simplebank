package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ly1999-hub/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := CreateJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotNil(t, maker)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	fmt.Println("CreateToken")
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	fmt.Println("verifyToken")
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssueAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomCurency(), time.Minute)
	require.NoError(t, err)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := CreateJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
