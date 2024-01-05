package fmhandler

import (
	"bytes"
	"encoding/json"
	"finance-manager-backend/cmd/finance-mngr/internal/application"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/jsonutils"
	"finance-manager-backend/cmd/finance-mngr/internal/repository/dbrepo"
	"finance-manager-backend/cmd/finance-mngr/internal/testingutils"
	"finance-manager-backend/cmd/finance-mngr/internal/validation"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var p testingutils.DockerTestPlatform
var fmh FinanceManagerHandler
var app application.Application

func TestMain(m *testing.M) {
	method := "validation_test.TestMain"
	fmlogger.Enter(method)

	//Setup testing platform using docker
	p = testingutils.Setup(m)
	db := &dbrepo.PostgresDBRepo{DB: p.DB}
	fmh = FinanceManagerHandler{
		DB:        db,
		JSONUtil:  &jsonutils.JSONUtil{},
		Validator: &validation.FinanceManagerValidator{DB: db},
		Auth:      testingutils.GetTestAuth(),
	}

	//Set application's handler
	app.Handler = &fmh
	app.Auth = fmh.Auth

	//Execute Code
	code := m.Run()

	//Tear down testing platform
	testingutils.TearDown(p)

	fmlogger.Exit(method)
	os.Exit(code)
}

func MakeRequest(method, url string, body interface{}, isAuthenticatedRequest bool, token string) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))

	if isAuthenticatedRequest {
		request.Header.Add("Authorization", "Bearer "+token)
	}

	writer := httptest.NewRecorder()
	app.Routes().ServeHTTP(writer, request)
	return writer

}
