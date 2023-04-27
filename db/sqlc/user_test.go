package db

import (
	"context"
	"testing"
	"time"

	"simpledice/util"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	owner := util.RandomOwner()

	user, err := testQueries.CreateUser(context.Background(), owner)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, owner, user.Username)
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
