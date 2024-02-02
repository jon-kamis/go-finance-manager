package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/authentication"
	"finance-manager-backend/internal/finance-mngr/models/restmodels"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// Authenticate godoc
// @title		Login
// @version 	1.0.0
// @Tags 		Authentication
// @Summary 	Login
// @Description Attempts to use passed credentials to authenticate with the application and generate JWT tokens
// @Accept		json
// @Produce 	json
// @Success 	200 {object} authentication.TokenPairs
// @Failure 	400 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/authenticate [post]
func (fmh *FinanceManagerHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload

	var requestPayload restmodels.LoginRequest

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

// RefreshToken godoc
// @title		Refresh Token
// @version 	1.0.0
// @Tags 		Authentication
// @Summary 	Refresh Token
// @Description Attempts to refresh Tokens using a refresh token
// @Accept		json
// @Produce 	json
// @Success 	200 {object} authentication.TokenPairs
// @Failure 	400 {object} jsonutils.JSONResponse
// @Failure 	401 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/refresh [get]
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

// Logout godoc
// @title		Logout
// @version 	1.0.0
// @Tags 		Authentication
// @Summary 	Logout
// @Description Returns an expired refresh cookie which prevents the user from re-authenticating
// @Accept		json
// @Produce 	json
// @Success 	200 
// @Router 		/logout [get]
func (fmh *FinanceManagerHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, fmh.Auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}
