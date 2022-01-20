package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "NewPassword"
	hashed, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	err = CheckPassword(password, hashed)
	require.NoError(t, err)

	passwordNew := "NewPassword2"
	err = CheckPassword(passwordNew, hashed)
	require.Error(t, err)

}
