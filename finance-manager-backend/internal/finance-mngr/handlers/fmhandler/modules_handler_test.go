package fmhandler

import (
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/internal/finance-mngr/testingutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIsModuleEnabled(t *testing.T) {
	method := "modules_handler_test.TestGetIsModuleEnabled"
	fmlogger.Enter(method)

	token := testingutils.GetAdminJWT(t)

	writer := MakeRequest(http.MethodGet, "/modules/stocks", nil, true, token)
	assert.Equal(t, http.StatusOK, writer.Code)

	var response models.ModuleEnabledResponse
	err := ReadResponse(writer, &response)
	assert.Nil(t, err)
	assert.False(t, response.Enabled)

	//Test not found
	writer = MakeRequest(http.MethodGet, "/modules/something-that-does-not-exist", nil, true, token)
	assert.Equal(t, http.StatusNotFound, writer.Code)

	fmlogger.Exit(method)
}