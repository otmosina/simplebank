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
	var a int64
	a = 1
	require.NotEmpty(t, a)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var limit int32 = 5
	var offset int32 = 0

	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().ListAccounts(gomock.Any(), db.ListAccountsParams{
		Limit:  limit,
		Offset: offset,
	}).
		Times(1).
		Return(accounts, nil)

	server := NewServer(store)

	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d", 1, 5)
	fmt.Println("==============================WE ARE LOOKING FOR THE NEXT URL==============================")
	fmt.Println(url)

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	server.router.ServeHTTP(recorder, req)
	// fmt.Println(err.Error())
	require.Equal(t, http.StatusOK, recorder.Code)
}
