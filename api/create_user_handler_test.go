package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	mockdb "github.com/otmosina/simplebank/db/mock"
	db "github.com/otmosina/simplebank/db/sqlc"
	"github.com/otmosina/simplebank/util"
	"github.com/stretchr/testify/require"
)

func randomUser() db.User {
	return db.User{
		Username:       util.RandomOwner(),
		HashedPassword: util.RandomString(6),
		Fullname:       util.RandomName(),
		Email:          util.RandomEmail(),
	}
}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}
func TestCreateUserAPI(t *testing.T) {
	//

	request := CreateUsersRequest{
		Username: util.RandomOwner(),
		Password: util.RandomString(6),
		Fullname: util.RandomName(),
		Email:    util.RandomEmail(),
	}

	badRequestNotAlphaNum := CreateUsersRequest{
		Username: "_________________",
		Password: util.RandomString(6),
		Fullname: util.RandomName(),
		Email:    util.RandomEmail(),
	}

	badRequestShortPass := CreateUsersRequest{
		Username: util.RandomOwner(),
		Password: util.RandomString(3),
		Fullname: util.RandomName(),
		Email:    util.RandomEmail(),
	}

	badRequestBadEmailTemplate := CreateUsersRequest{
		Username: util.RandomOwner(),
		Password: util.RandomString(6),
		Fullname: util.RandomName(),
		Email:    "IAMNOTEMAIL",
	}
	user := randomUser()
	testCases := []testCaseUserCreate{
		{
			name:    "OK",
			request: request,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(db.CreateUserParams{
					Username:       request.Username,
					HashedPassword: request.Password,
					Fullname:       request.Fullname,
					Email:          request.Email,
				}, request.Password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:    "UsernameIsNotAlphaNum",
			request: badRequestNotAlphaNum,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "PasswordTooShort",
			request: badRequestShortPass,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "NotEmailInRequest",
			request: badRequestBadEmailTemplate,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "InternalServerError",
			request: request,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(db.CreateUserParams{
					Username:       request.Username,
					HashedPassword: request.Password,
					Fullname:       request.Fullname,
					Email:          request.Email,
				}, request.Password)).
					Times(1).
					Return(db.User{}, pq.ErrSSLNotSupported)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	// var url string
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			store := getMockStore(t)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()
			// url = "/accounts"
			url := fmt.Sprintf("/users")
			bytesRequest, _ := json.Marshal(tc.request)
			reader := bytes.NewReader(bytesRequest)
			// reader := strings.NewReader(bytesRequest)

			request, err := http.NewRequest(http.MethodPost, url, reader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}

}
