package db

import (
	"context"
	"testing"
	"time"

	"github.com/otmosina/simplebank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		Fullname:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	// if err != nil {
	// 	log.Fatal("Error when createAccoutn")
	// }
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username, "they should be eq")
	require.Equal(t, arg.HashedPassword, user.HashedPassword, "they should be eq")
	require.Equal(t, arg.Fullname, user.Fullname, "they should be eq")
	require.Equal(t, arg.Email, user.Email, "they should be eq")

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())
	// require.NotZero(t, user.PasswordChangedAt)
	// assert.Equal(t, user.ID, int64(2), "they should be eq")
	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)

	require.Equal(t, user1.Username, user2.Username, "they should be eq")
	require.Equal(t, user1.HashedPassword, user2.HashedPassword, "they should be eq")
	require.Equal(t, user1.Fullname, user2.Fullname, "they should be eq")
	require.Equal(t, user1.Email, user2.Email, "they should be eq")

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}
