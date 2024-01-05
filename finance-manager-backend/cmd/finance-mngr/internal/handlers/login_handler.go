package handlers

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/authentication"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func (fmh *FinanceManagerHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload
	var requestPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := fmh.JSONUtil.ReadJSON(w, r, &requestPayload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate user against database
	user, err := fmh.DB.GetUserByUsername((requestPayload.Username))
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		fmh.JSONUtil.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
	}

	// load user roles
	roles, err := fmh.DB.GetAllUserRoles(user.ID)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err)
		return
	}

	// Convert user role codes into csv string
	var roleSlice []string
	for _, role := range roles {
		roleSlice = append(roleSlice, role.Code)
	}

	roleStr := strings.Join(roleSlice, ",")

	// create a jwt user
	u := authentication.JwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Roles:     roleStr,
	}

	//generate tokens
	tokens, err := fmh.Auth.GenerateTokenPair(&u)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err)
		return
	}

	refreshCookie := fmh.Auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, tokens)
}

func (fmh *FinanceManagerHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == fmh.Auth.CookieName {
			claims := &authentication.Claims{}
			refreshToken := cookie.Value

			//parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(fmh.Auth.Secret), nil
			})

			if err != nil {
				fmh.JSONUtil.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
			}

			//get the user id from the token claims
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				fmh.JSONUtil.ErrorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
			}

			user, err := fmh.DB.GetUserByID(userID)
			if err != nil {
				fmh.JSONUtil.ErrorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
			}

			u := authentication.JwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := fmh.Auth.GenerateTokenPair(&u)
			if err != nil {
				fmh.JSONUtil.ErrorJSON(w, errors.New("error generating tokens"), http.StatusInternalServerError)
			}

			http.SetCookie(w, fmh.Auth.GetRefreshCookie(tokenPairs.RefreshToken))

			fmh.JSONUtil.WriteJSON(w, http.StatusOK, tokenPairs)
		}
	}
}

func (fmh *FinanceManagerHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, fmh.Auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}
