package api

import (
	"bytes"
	"database/sql"
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

func TestCreateAccoutAPI(t *testing.T) {
	//

	request := CreateAccountsRequest{
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
	}
	account := randomAccount()

	testCases := []testCaseCreate{
		{
			name:    "OK",
			request: request,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{
					Owner:    request.Owner,
					Currency: request.Currency,
					Balance:  0,
				})).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name:    "BadRequest",
			request: CreateAccountsRequest{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:    "StatusInternalServerError2",
			request: request,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{
					Owner:    request.Owner,
					Currency: request.Currency,
					Balance:  0,
				})).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		// {
		// 	name:      "GetAccountError",
		// 	accountID: account.ID,
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).
		// 			Times(1).
		// 			Return(db.Account{}, sql.ErrNoRows)
		// 		store.EXPECT().DeleteAccount(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusNotFound, recorder.Code)
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
			url := fmt.Sprintf("/accounts")
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
