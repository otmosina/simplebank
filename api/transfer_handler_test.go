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

func TestTransferRequestAPI(t *testing.T) {
	//

	// type transferRequest struct {
	// 	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	// 	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	// 	Amount        int64  `json:"amount" binding:"required,gt=0"`
	// 	Currency      string `json:"currency" binding:"required,oneof=USD,RUB,IDR"`
	// }

	currency := util.RandomCurrency()

	account1 := randomAccount()
	account2 := randomAccount()

	account1.Currency = currency
	account2.Currency = currency

	minBalance := account1.Balance
	if minBalance > account2.Balance {
		minBalance = account2.Balance
	}
	amount := util.RandomInt(1, minBalance)

	postRequest := TransferParamsRequest{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
		Currency:      "USD",
	}

	transferParams := db.TransferTxParams{
		FromAccountID: postRequest.FromAccountID,
		ToAccountID:   postRequest.ToAccountID,
		Amount:        postRequest.Amount,
	}

	testCases := []testCaseTransfer{
		{
			name:    "OK",
			request: postRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account1.ID).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccount(gomock.Any(), account2.ID).
					Times(1).
					Return(account2, nil)

				store.EXPECT().TransferTx(gomock.Any(), transferParams).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Wrong currency",
			request: TransferParamsRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      "WRONG",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	// var url string = "/transfers"
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			store := getMockStore(t)
			tc.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			bytesRequest, _ := json.Marshal(tc.request)
			reader := bytes.NewReader(bytesRequest)

			url := fmt.Sprintf("/transfers")
			// fmt.Println()
			request, err := http.NewRequest(http.MethodPost, url, reader)
			// fmt.Println(request.Body)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}
}
