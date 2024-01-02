package handlers

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/authentication"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func (fmh *FinanceManagerHandler) Register(w http.ResponseWriter, r *http.Request) {
	method := "login_handler.Register"
	fmt.Printf("[ENTER %s]\n", method)

	var payload models.User

	// Read in user from payload
	err := fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err)
		fmt.Printf("[%s] failed to read JSON payload", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	// Validate if user can be entered
	err = fmh.Validator.IsValidToEnterNewUser(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err)
		fmt.Printf("[%s] %s", method, err)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	// Retrieve default role of 'user'
	role, err := fmh.DB.GetRoleByCode("user")
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err)
		fmt.Printf("[%s] failed to load default role for new user\n", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	userId, err := fmh.DB.InsertUser(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err)
		fmt.Printf("[%s] failed to insert new user", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	userRole := models.UserRole{
		UserId:       userId,
		RoleId:       role.ID,
		Code:         role.Code,
		CreateDt:     time.Now(),
		LastUpdateDt: time.Now(),
	}

	_, err = fmh.DB.InsertUserRole(userRole)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err)
		fmt.Printf("[%s] failed to insert new userRole", method)
		fmt.Printf("[EXIT %s]\n", method)
		return
	}

	fmt.Printf("[EXIT %s]\n", method)
	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, "new user was created successfully")
}
