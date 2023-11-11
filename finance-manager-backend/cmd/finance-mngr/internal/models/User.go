package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"user_name"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	CreateDt     time.Time `json:"-"`
	LastUpdateDt time.Time `json:"-"`
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			//invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
