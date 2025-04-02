package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testPassword = "mySecurePassword123"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword(testPassword)

	assert.NoError(t, err)

	assert.NotEmpty(t, hash)

	assert.NotEqual(t, testPassword, hash)
}

func TestCheckPasswordHash_Valid(t *testing.T) {
	hash, err := HashPassword(testPassword)
	assert.NoError(t, err)

	isValid := CheckPasswordHash(testPassword, hash)
	assert.True(t, isValid)
}

func TestCheckPasswordHash_Invalid(t *testing.T) {
	password := "mySecurePassword123"
	incorrectPassword := "incorrectPassword"

	hash, err := HashPassword(password)
	assert.NoError(t, err)

	isValid := CheckPasswordHash(incorrectPassword, hash)
	assert.False(t, isValid)
}
