package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
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
				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(db.CreateUserParams{
					Username:       request.Username,
					HashedPassword: request.Password,
					Fullname:       request.Fullname,
					Email:          request.Email,
				})).
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

		// {
		// 	name:    "BadRequest",
		// 	request: CreateAccountsRequest{},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
		// {
		// 	name:    "StatusInternalServerError2",
		// 	request: request,
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{
		// 			Owner:    request.Owner,
		// 			Currency: request.Currency,
		// 			Balance:  0,
		// 		})).
		// 			Times(1).
		// 			Return(db.Account{}, sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		// 	},
		// },
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
