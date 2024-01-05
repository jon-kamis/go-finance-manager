package fmhandler

import (
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"net/http"
	"time"
)

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
