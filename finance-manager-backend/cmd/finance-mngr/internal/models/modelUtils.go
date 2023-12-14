package models

import (
	"errors"
	"fmt"
)

func returnError(msg string, method string) error {
	fmt.Printf("[%s] %s\n", method, msg)
	fmt.Printf("[EXIT %s]\n", method)
	return errors.New(msg)
}
