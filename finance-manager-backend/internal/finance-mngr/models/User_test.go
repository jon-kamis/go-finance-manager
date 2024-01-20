package models

import (
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordMatches(t *testing.T) {
	method := "User_test.TestPasswordMatches"
	fmlogger.Enter(method)

	password := "abc123"
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	u := User{
		Password: string(encryptedPass),
	}

	//Good password
	success, err := u.PasswordMatches(password)
	assert.True(t, success)
	assert.Nil(t, err)

	//Bad Password
	success, err = u.PasswordMatches("wrongpassword")
	assert.False(t, success)
	assert.Nil(t, err)

	fmlogger.Exit(method)
}
