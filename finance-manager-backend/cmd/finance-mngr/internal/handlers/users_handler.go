package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (fmh *FinanceManagerHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	method := "handlers_user.GetUserByID"
	fmt.Printf("[ENTER %s]\n", method)

	idStr := chi.URLParam(r, "userId")
	var id int

	loggedInUserId, err := fmh.Auth.GetLoggedInUserId(w, r)

	if err != nil {
		fmt.Printf("[%v] unexpected error occured when fetching logged in user: %v\n", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching user"), http.StatusInternalServerError)
		return
	}

	if idStr == "me" {
		id = loggedInUserId
	} else {

		id, err = strconv.Atoi(idStr)
		if err != nil {
			fmt.Printf("[%v] the supplied id was invalid and returned the error: %v\n", method, err)
			fmt.Printf("[EXIT %s]\n", method)
			fmh.JSONUtil.ErrorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
			return
		}

		if id != loggedInUserId {
			isValid, err := fmh.Validator.IsValidToViewOtherUserData(loggedInUserId)

			if err != nil {
				fmt.Printf("[%v] unexpected error occured when fetching user: %v\n", method, err)
				fmt.Printf("[EXIT %s]\n", method)
				fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching user"), http.StatusInternalServerError)
				return
			} else if !isValid {
				fmt.Printf("[%v] user is forbidden from viewing other user data\n", method)
				fmt.Printf("[EXIT %s]\n", method)
				fmh.JSONUtil.ErrorJSON(w, errors.New("forbidden"), http.StatusForbidden)
				return
			}

		}
	}

	user, err := fmh.DB.GetUserByID(id)

	if err != nil {
		fmt.Printf("[%v] the requested user was not found\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, errors.New("user not found"), http.StatusNotFound)
		return
	}

	fmh.JSONUtil.WriteJSON(w, http.StatusOK, user)
}

func (fmh *FinanceManagerHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	method := "handlers_users.GetAllUsers"
	fmt.Printf("[ENTER %s]\n", method)

	search := r.URL.Query().Get("search")
	users, err := fmh.DB.GetAllUsers(search)

	if err != nil {
		fmt.Printf("[%s] caught unexpected error when attempting to fetch users from database\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching users"), http.StatusInternalServerError)
		return
	}

	fmt.Printf("[EXIT %s]", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, users)
}
