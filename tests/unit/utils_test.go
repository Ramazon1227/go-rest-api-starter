package unit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ramazon1227/go-rest-api-starter/pkg/utils"
)

func TestHashPassword(t *testing.T) {
	hash, err := utils.HashPassword("secret123")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "secret123", hash)
}

func TestHashPassword_DifferentEachCall(t *testing.T) {
	h1, err := utils.HashPassword("same")
	require.NoError(t, err)
	h2, err := utils.HashPassword("same")
	require.NoError(t, err)
	// bcrypt salts each hash, so they must differ
	assert.NotEqual(t, h1, h2)
}

func TestCheckPassword_Correct(t *testing.T) {
	hash, err := utils.HashPassword("mypassword")
	require.NoError(t, err)
	assert.True(t, utils.CheckPassword(hash, "mypassword"))
}

func TestCheckPassword_Wrong(t *testing.T) {
	hash, err := utils.HashPassword("mypassword")
	require.NoError(t, err)
	assert.False(t, utils.CheckPassword(hash, "wrongpassword"))
}

func TestCheckPassword_EmptyPassword(t *testing.T) {
	hash, err := utils.HashPassword("nonempty")
	require.NoError(t, err)
	assert.False(t, utils.CheckPassword(hash, ""))
}

func TestGenerateRandomPassword_Length(t *testing.T) {
	for _, length := range []int{6, 8, 12, 16} {
		pass, err := utils.GenerateRandomPassword(length)
		require.NoError(t, err)
		assert.Len(t, pass, length)
	}
}

func TestGenerateRandomPassword_Unique(t *testing.T) {
	p1, err := utils.GenerateRandomPassword(12)
	require.NoError(t, err)
	p2, err := utils.GenerateRandomPassword(12)
	require.NoError(t, err)
	assert.NotEqual(t, p1, p2)
}
