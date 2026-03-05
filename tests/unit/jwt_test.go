package unit_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ramazon1227/go-rest-api-starter/models"
	"github.com/Ramazon1227/go-rest-api-starter/pkg/jwt"
)

func init() {
	jwt.SigningKey = []byte("test-secret-key")
}

func testUser() *models.User {
	return &models.User{
		Id:    "user-123",
		Email: "test@example.com",
		Role:  "SYSTEM_ADMIN",
	}
}

func TestGenerateToken(t *testing.T) {
	token, err := jwt.GenerateToken(testUser())
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken_Valid(t *testing.T) {
	token, err := jwt.GenerateToken(testUser())
	require.NoError(t, err)

	parsed, err := jwt.ValidateToken(token)
	require.NoError(t, err)
	assert.True(t, parsed.Valid)
}

func TestValidateToken_Tampered(t *testing.T) {
	token, err := jwt.GenerateToken(testUser())
	require.NoError(t, err)

	_, err = jwt.ValidateToken(token + "tampered")
	assert.Error(t, err)
}

func TestValidateToken_WrongKey(t *testing.T) {
	token, err := jwt.GenerateToken(testUser())
	require.NoError(t, err)

	// Temporarily swap signing key
	jwt.SigningKey = []byte("wrong-key")
	_, err = jwt.ValidateToken(token)
	assert.Error(t, err)

	// Restore
	jwt.SigningKey = []byte("test-secret-key")
}

func TestExtractClaims(t *testing.T) {
	user := testUser()
	token, err := jwt.GenerateToken(user)
	require.NoError(t, err)

	claims, err := jwt.ExtractClaims(token)
	require.NoError(t, err)

	assert.Equal(t, user.Id, claims["user_id"])
	assert.Equal(t, user.Email, claims["email"])
	assert.Equal(t, user.Role, claims["role"])
}

func TestExtractClaims_ExpiredToken(t *testing.T) {
	// Temporarily set a very short expiry
	original := jwt.ExpiryDuration
	jwt.ExpiryDuration = -time.Second
	defer func() { jwt.ExpiryDuration = original }()

	token, err := jwt.GenerateToken(testUser())
	require.NoError(t, err)

	_, err = jwt.ExtractClaims(token)
	assert.Error(t, err)
}

func TestGetTokenExpiryTime(t *testing.T) {
	before := time.Now()
	expiry := jwt.GetTokenExpiryTime()
	after := time.Now()

	assert.True(t, expiry.After(before))
	assert.True(t, expiry.After(after))
}

func TestTokenBlacklist_NotBlacklisted(t *testing.T) {
	assert.False(t, jwt.IsTokenBlacklisted("some-fresh-token"))
}

func TestTokenBlacklist_AfterInvalidate(t *testing.T) {
	token, err := jwt.GenerateToken(testUser())
	require.NoError(t, err)

	err = jwt.InvalidateToken(token)
	require.NoError(t, err)

	assert.True(t, jwt.IsTokenBlacklisted(token))
}
