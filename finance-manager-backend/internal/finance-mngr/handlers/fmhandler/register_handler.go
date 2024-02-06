package fmhandler

import (
	"finance-manager-backend/internal/finance-mngr/constants"
	"finance-manager-backend/internal/finance-mngr/models"
	"net/http"
	"time"

	"github.com/jon-kamis/klogger"
)

// Register godoc
// @title		Register
// @version 	1.0.0
// @Tags 		Authentication
// @Summary 	Register
// @Description Attempts to register a new user into the application
// @Param		user body models.User true "The User to Register"
// @Accept		json
// @Produce 	json
// @Success 	200 {object} jsonutils.JSONResponse
// @Failure 	400 {object} jsonutils.JSONResponse
// @Failure 	500 {object} jsonutils.JSONResponse
// @Router 		/register [post]
func (fmh *FinanceManagerHandler) Register(w http.ResponseWriter, r *http.Request) {
	method := "register_handler.Register"
	klogger.Enter(method)

	var payload models.User

	// Read in user from payload
	err := fmh.JSONUtil.ReadJSON(w, r, &payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err)
		klogger.ExitError(method, constants.FailedToParseJsonBodyError, err)
		return
	}

	// Validate if user can be entered
	err = fmh.Validator.IsValidToEnterNewUser(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err)
		klogger.ExitError(method, constants.GenericUnprocessableEntityErrLog, err)
		return
	}

	// Retrieve default role of 'user'
	role, err := fmh.DB.GetRoleByCode("user")
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err)
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return
	}

	userId, err := fmh.DB.InsertUser(payload)
	if err != nil {
		fmh.JSONUtil.ErrorJSON(w, err)
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
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
		klogger.ExitError(method, constants.UnexpectedSQLError, err)
		return
	}

	klogger.Exit(method)
	fmh.JSONUtil.WriteJSON(w, http.StatusAccepted, constants.SuccessMessage)
}
