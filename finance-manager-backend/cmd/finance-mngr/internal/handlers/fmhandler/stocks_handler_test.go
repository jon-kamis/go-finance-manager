package fmhandler

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"finance-manager-backend/cmd/finance-mngr/internal/testingutils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIsStocksEnabled(t *testing.T) {
	method := "stocks_handler_test.TestGetIsStocksEnabled"
	fmlogger.Enter(method)

	token := testingutils.GetAdminJWT(t)

	writer := MakeRequest(http.MethodGet, "/modules/stocks", nil, true, token)
	assert.Equal(t, http.StatusOK, writer.Code)

	var response models.StocksEnabledResponse
	err := ReadResponse(writer, &response)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	assert.False(t, response.Enabled)

	fmlogger.Exit(method)
}
