package api

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/otmosina/simplebank/db/mock"
	db "github.com/otmosina/simplebank/db/sqlc"
)

// testCaseBase {
// 	name          string
// 	buildStubs    func(store *mockdb.MockStore)
// 	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// }

type testCase struct {
	name          string
	accountID     int64
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

type testCaseCreate struct {
	name          string
	request       CreateAccountsRequest
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

type testCaseIndex struct {
	name          string
	pageID        int32
	pageSize      int32
	request       db.ListAccountsParams
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

func TestMain(m *testing.M) {
	fmt.Println("MAIN_TEST API API API API API API API API API API API API API API API API ")
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func getMockStore(t *testing.T) *mockdb.MockStore {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	return mockdb.NewMockStore(ctrl)

}
