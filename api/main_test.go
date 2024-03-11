package api

import (
	"testing"
	"time"

	db "github.com/ly1999-hub/simplebank/sqlc"
	"github.com/ly1999-hub/simplebank/util"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	if err != nil {
		return nil
	}
	return server
}
