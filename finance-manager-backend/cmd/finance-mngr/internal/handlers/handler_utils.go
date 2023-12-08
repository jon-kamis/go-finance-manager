package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (fmh *FinanceManagerHandler) GetAndValidateUserId(idStr string, canViewOtherUserData bool, w http.ResponseWriter, r *http.Request) (int, error) {
	method := "handler_utils.getAndvalidateUserId"
	fmt.Printf("[ENTER %s]\n", method)

	// Get loggedIn userId
	loggedInUserId, err := fmh.Auth.GetLoggedInUserId(w, r)
	var id int

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when retrieving logged in userId from HTTP Header: %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return -1, errors.New("failed to retrieve logged in userId")
	}

	if idStr == "me" {
		id = loggedInUserId
	} else {

		id, err = strconv.Atoi(idStr)
		if err != nil {

			return -1, errors.New("invalid id")
		}

		if id != loggedInUserId && !canViewOtherUserData {
			fmt.Printf("[%s] user is forbidden from viewing other user's data\n", method)
			fmt.Printf("[EXIT %s]\n", method)
			return -1, errors.New("forbidden")
		}
	}

	fmt.Printf("[EXIT %s]\n", method)
	return id, nil
}
