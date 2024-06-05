package pkg

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestVerifyToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	email := "test@example.com"
	tokenString, err := CreateToken(email)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	returnedEmail, err := VerifyToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, email, returnedEmail)

	_, err = VerifyToken("invalidToken")
	assert.Error(t, err)
	assert.Equal(t, "token is invalid", err.Error())

	_, err = VerifyToken("")
	assert.Error(t, err)
	assert.Equal(t, "token is empty", err.Error())

	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(-time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredTokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	assert.NoError(t, err)
	_, err = VerifyToken(expiredTokenString)
	assert.Error(t, err)
	assert.Equal(t, "token is invalid", err.Error())
}
