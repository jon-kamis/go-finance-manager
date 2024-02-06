package models

import (
	"errors"
	"time"

	"github.com/jon-kamis/klogger"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	CreateDt     time.Time `json:"-"`
	LastUpdateDt time.Time `json:"-"`
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	method := "User.PasswordMatches"
	klogger.Enter(method)

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			klogger.ExitError(method, "invalid password")
			return false, nil
		default:
			klogger.ExitError(method, "unexpected error")
			return false, err
		}
	}

	klogger.Exit(method)
	return true, nil
}
