package models

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"time"

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
	fmlogger.Enter(method)

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			fmlogger.ExitError(method, "invalid password", err)
			return false, nil
		default:
			fmlogger.ExitError(method, "unexpected error", err)
			return false, err
		}
	}

	fmlogger.Exit(method)
	return true, nil
}
