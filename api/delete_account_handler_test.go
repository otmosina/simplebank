package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestDeleteAccoutAPI(t *testing.T) {
	//
	var accountID int64
	account := randomAccount()
	accountID = account.ID

	store := getMockStore(t) // mockdb.NewMockStore(ctrl)
	// buildStubs
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountID)).
		Times(1).
		Return(account, nil)

	store.EXPECT().DeleteAccount(gomock.Any(), gomock.Eq(accountID)).
		Times(1).
		Return(nil)
	// buildStubs

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/account/%d", accountID)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

}
