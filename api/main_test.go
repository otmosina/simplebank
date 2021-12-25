package api

import (
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/otmosina/simplebank/db/mock"
)

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
