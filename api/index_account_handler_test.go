package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/otmosina/simplebank/db/mock"
	db "github.com/otmosina/simplebank/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestIndexAccountAPI(t *testing.T) {

	account1 := randomAccount()
	account2 := randomAccount()
	var accounts []db.Account
	accounts = append(accounts, account1)
	accounts = append(accounts, account2)

	var limit, offset, pageID, pageSize int32 // = 5
	//var  offset int32// = 0

	pageID = 1
	pageSize = 5
	limit = pageSize
	offset = (pageID - 1) * pageSize

	testCases := []testCaseIndex{
		{
			name: "OK",
			request: db.ListAccountsParams{
				Limit:  limit,
				Offset: offset,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), db.ListAccountsParams{
					Limit:  limit,
					Offset: offset,
				}).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d", 1, 5)
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}

	// fmt.Println(err.Error())
	// require.Equal(t, http.StatusOK, recorder.Code)
}
