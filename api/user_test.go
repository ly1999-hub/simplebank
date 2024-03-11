package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/ly1999-hub/simplebank/db/mock"
	db "github.com/ly1999-hub/simplebank/sqlc"
	"github.com/ly1999-hub/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAPI(t *testing.T) {
	user, _ := randomUser(t)

	testCase := []struct {
		name         string
		body         gin.H
		buildStubs   func(store *mockdb.MockStore)
		checkRespone func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  "huuly1999",
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkRespone: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.User{}, nil)
			},
			checkRespone: func(recorder *httptest.ResponseRecorder) {
				fmt.Println("recorder:", recorder)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCase {

		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			//build stubs
			tc.buildStubs(store)
			//start test server send Request
			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			url := "/users"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkRespone(recorder)
		})
	}

}

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	hashedpassword, err := util.HashPassword(password)
	require.NoError(t, err)
	user := db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedpassword,
		FullName:       util.RandomOwner(),
		Email:          util.CreateRandomEmail(),
	}
	return user, hashedpassword

}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var gotAccount db.User

	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, gotAccount, user)
}
