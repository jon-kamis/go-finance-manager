package handlers

import (
	"errors"
	"finance-manager-backend/cmd/finance-mngr/internal/constants"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (fmh *FinanceManagerHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	method := "user_role_handler.GetUserRoles"
	fmlogger.Enter(method)

	//Read ID from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching user roles", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	userRoles, err := fmh.DB.GetAllUserRoles(userId)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching user roles", err)
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching user roles"), http.StatusInternalServerError)
		return
	}

	fmlogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, userRoles)
}

func (fmh *FinanceManagerHandler) AddUserRoles(w http.ResponseWriter, r *http.Request) {
	method := "user_role_handler.AddUserRoles"
	fmlogger.Enter(method)

	//Read IDs from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	roleId := chi.URLParam(r, "roleId")

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured validating userId", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	//Fetch role
	role, err := fmh.DB.GetRoleById(roleId)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching role", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Check if the user already has the requested role
	hasRole, err := fmh.Validator.CheckIfUserHasRole(userId, role.Code)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when checking if user has requested Role", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if hasRole {
		err = errors.New("user already has requested role")
		fmlogger.ExitError(method, "user already has requested role", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
	}

	// Fetch user
	user, err := fmh.DB.GetUserByID(userId)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching user", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Insert new user role
	newRole := models.UserRole{
		UserId: user.ID,
		RoleId: role.ID,
		Code:   role.Code,
	}

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when fetching user", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_, err = fmh.DB.InsertUserRole(newRole)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when inserting user role", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	fmlogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, "success")
}

func (fmh *FinanceManagerHandler) DeleteUserRoles(w http.ResponseWriter, r *http.Request) {
	method := "user_role_handler.DeleteUserRoles"
	fmlogger.Enter(method)

	//Read IDs from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured validating userId", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	roleId, err := strconv.Atoi(chi.URLParam(r, "roleId"))

	if err != nil {

		fmlogger.ExitError(method, "invalid roleId", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Check if the user role exists and belongs to the given user
	err = fmh.Validator.UserRoleExistsAndBelongsToUser(int(roleId), userId)

	if err != nil {
		fmlogger.ExitError(method, "user role does not exist or does not belong to given user", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		return
	}

	userRole, err := fmh.DB.GetUserRoleByRoleIDAndUserID(roleId, userId)
	
	if err != nil {
		fmlogger.ExitError(method, "error during db call", err)
		err = errors.New(constants.UnexpectedSQLError)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		return
	}

	err = fmh.DB.DeleteUserRoleByID(userRole.ID)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error occured when inserting user role", err)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	fmlogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, "success")
}
