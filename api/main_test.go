package api

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/otmosina/simplebank/db/mock"
)

type testCase struct {
	name          string
	accountID     int64
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
