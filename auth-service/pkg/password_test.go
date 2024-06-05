package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password"
	hash, _ := HashPassword(password)
	assert.True(t, CheckPasswordHash(password, hash))
}
