package fmhandler

import (
	"errors"
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jon-kamis/klogger"
)

// GetUserRoles godoc
// @title		Get All User Roles
// @version 	1.0.0
// @Tags 		User Roles
// @Summary 	Get All User Roles
// @Description Returns an array of UserRole objects belonging to a given user
// @Param		userId path int true "ID of the user we are searching for"
// @Param		search query string false "Search for user role by role name"
// @Produce 	json
// @Success 	200 {array} models.UserRole
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/roles [get]
func (fmh *FinanceManagerHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	method := "user_role_handler.GetUserRoles"
	klogger.Enter(method)

	//Read ID from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	userRoles, err := fmh.DB.GetAllUserRoles(userId)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, errors.New("unexpected error occured when fetching user roles"), http.StatusInternalServerError)
		klogger.ExitError(method, constants.FailedToRetrieveEntityError, err)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, userRoles)
}

// AddUserRoles godoc
// @title		Add User Role
// @version 	1.0.0
// @Tags 		User Roles
// @Summary 	Add User Role
// @Description Adds a new role to a User
// @Param		userId path int true "ID of the user to add a role to"
// @Param		roleId path int true "ID of the role to add to the user"
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/roles/{roleId} [post]
func (fmh *FinanceManagerHandler) AddUserRoles(w http.ResponseWriter, r *http.Request) {
	method := "user_role_handler.AddUserRoles"
	klogger.Enter(method)

	//Read IDs from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)
	roleId := chi.URLParam(r, "roleId")

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	//Fetch role
	role, err := fmh.DB.GetRoleById(roleId)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.FailedToRetrieveEntityError, err)
		return
	}

	// Check if the user already has the requested role
	hasRole, err := fmh.Validator.CheckIfUserHasRole(userId, role.Code)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.GenericUnexpectedErrorLog, err)
		return
	}

	if hasRole {
		err = errors.New("user already has requested role")
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		klogger.ExitError(method, "user already has requested role")
		return
	}

	// Fetch user
	user, err := fmh.DB.GetUserByID(userId)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.FailedToRetrieveEntityError, err)
		return
	}

	// Insert new user role
	newRole := models.UserRole{
		UserId: user.ID,
		RoleId: role.ID,
		Code:   role.Code,
	}

	_, err = fmh.DB.InsertUserRole(newRole)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.FailedToSaveEntityError, err)
		return
	}

	fmh.JSONUtil.WriteJSON(w, http.StatusOK, constants.SuccessMessage)
	klogger.Exit(method)
}

// DeleteUserRoles godoc
// @title		Remove User Role
// @version 	1.0.0
// @Tags 		User Roles
// @Summary 	Remove User Role
// @Description Removes a role from a a User
// @Param		userId path int true "ID of the user to remove a role from"
// @Param		roleId path int true "ID of the role to remove from the user"
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	403 {object} jsonutils.JSONResponse
// @Failure 	404 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/users/{userId}/roles/{roleId} [delete]
func (fmh *FinanceManagerHandler) DeleteUserRoles(w http.ResponseWriter, r *http.Request) {
	method := "user_role_handler.DeleteUserRoles"
	klogger.Enter(method)

	//Read IDs from url
	userId, err := fmh.GetAndValidateUserId(chi.URLParam(r, "userId"), w, r)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.EntityDoesNotBelongToUserError, err)
		return
	}

	roleId, err := strconv.Atoi(chi.URLParam(r, "roleId"))

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusBadRequest)
		klogger.ExitError(method, "invalid roleId")
		return
	}

	// Check if the user role exists and belongs to the given user
	err = fmh.Validator.UserRoleExistsAndBelongsToUser(int(roleId), userId)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		klogger.ExitError(method, "user role does not exist or does not belong to given user")
		return
	}

	userRole, err := fmh.DB.GetUserRoleByRoleIDAndUserID(roleId, userId)

	if err != nil {
		err = errors.New(constants.UnexpectedSQLError)
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusUnprocessableEntity)
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return
	}

	err = fmh.DB.DeleteUserRoleByID(userRole.ID)

	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err, http.StatusInternalServerError)
		klogger.ExitError(method, constants.FailedToDeleteEntityError, err)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusOK, "success")
}
