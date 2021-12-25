package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/otmosina/simplebank/db/mock"
	db "github.com/otmosina/simplebank/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestDeleteAccoutAPI(t *testing.T) {
	//
	var accountID int64
	account := randomAccount()
	accountID = account.ID

	testCases := []testCase{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountID)).
					Times(1).
					Return(account, nil)

				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Eq(accountID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name:      "BadRequest",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "GetAccountError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "DeleteAccountError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(account, nil)
				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			store := getMockStore(t)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/account/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodPost, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}

}
