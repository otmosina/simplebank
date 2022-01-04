package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	db "github.com/otmosina/simplebank/db/sqlc"
	"github.com/otmosina/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestTransferRequestAPI(t *testing.T) {
	//

	type transferRequest struct {
		FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
		ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
		Amount        int64  `json:"amount" binding:"required,gt=0"`
		Currency      string `json:"currency" binding:"required,oneof=USD,RUB,IDR"`
	}

	currency := util.RandomCurrency()

	account1 := randomAccount()
	account2 := randomAccount()

	account1.Currency = currency
	account2.Currency = currency

	minBalance := account1.Balance
	if minBalance > account2.Balance {
		minBalance = account2.Balance
	}
	amount := util.RandomInt(0, minBalance)

	store := getMockStore(t)

	store.EXPECT().
		GetAccount(gomock.Any(), account1.ID).
		Times(1).
		Return(account1, nil)

	store.EXPECT().
		GetAccount(gomock.Any(), account2.ID).
		Times(1).
		Return(account2, nil)

	transferParams := db.TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	}

	postRequest := transferRequest{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
		Currency:      "USD",
	}

	store.EXPECT().TransferTx(gomock.Any(), transferParams).
		Times(1)

	server := NewServer(store)
	recorder := httptest.NewRecorder()
	var url = "/transfers"
	// url := fmt.Sprintf("/transfers")
	bytesRequest, _ := json.Marshal(postRequest)
	reader := bytes.NewReader(bytesRequest)
	request, err := http.NewRequest(http.MethodPost, url, reader)
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

}
