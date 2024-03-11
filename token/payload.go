package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrExpireToken = errors.New("token has expire")
var ErrInvalidToken = errors.New("invalid Token")
var ErrInvalidLenSecretKey = errors.New("invalid key size")

// payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssueAt   time.Time `json:"issue_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	fmt.Println("new idToken")
	idToken, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("error newRandom uuid")
		return nil, err
	}
	payload := &Payload{
		ID:        idToken,
		Username:  username,
		IssueAt:   time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	fmt.Println(payload)
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt.UTC()) {
		return ErrExpireToken
	}
	return nil
}
