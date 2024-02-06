package fmhandler

import (
	"bytes"
	"encoding/json"
	"finance-manager-backend/internal/finance-mngr/application"
	"finance-manager-backend/internal/finance-mngr/jsonutils"
	"finance-manager-backend/internal/finance-mngr/repository/dbrepo"
	"finance-manager-backend/internal/finance-mngr/service/fmservice"
	"finance-manager-backend/internal/finance-mngr/service/polygonservice"
	"finance-manager-backend/internal/finance-mngr/validation"
	"finance-manager-backend/test"
	"finance-manager-backend/test/logtest"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jon-kamis/klogger"
)

var p test.DockerTestPlatform
var fmh FinanceManagerHandler
var app application.Application

func TestMain(m *testing.M) {
	logtest.SetKloggerTestFileNameEnv()

	method := "validation_test.TestMain"
	klogger.Enter(method)

	//Setup testing platform using docker
	p = test.Setup(m)
	db := &dbrepo.PostgresDBRepo{DB: p.DB}
	fmh = FinanceManagerHandler{
		DB:              db,
		JSONUtil:        &jsonutils.JSONUtil{},
		Validator:       &validation.FinanceManagerValidator{DB: db},
		Auth:            test.GetTestAuth(),
		Service:         &fmservice.FMService{},
		ExternalService: &polygonservice.PolygonService{},
	}

	//Set application's handler
	app.Handler = &fmh
	app.Auth = fmh.Auth

	//Execute Code
	code := m.Run()

	//Tear down testing platform
	test.TearDown(p)

	klogger.Exit(method)
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

func ReadResponse(w *httptest.ResponseRecorder, j interface{}) error {
	method := "FinanceManagerHandler_test.ReadResponse"
	klogger.Enter(method)

	err := json.NewDecoder(w.Body).Decode(&j)

	if err != nil {
		klogger.ExitError(method, "unexpected error", err)
		return err
	}

	klogger.Exit(method)
	return nil
}
