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
	"github.com/stretchr/testify/assert"
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
		DB:            db,
		JSONUtil:      &jsonutils.JSONUtil{},
		Validator:     &validation.FinanceManagerValidator{DB: db},
		Auth:          testingutils.GetTestAuth(),
		StocksEnabled: false,
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

func ReadResponse(w *httptest.ResponseRecorder, j interface{}) error{
	method := "FinanceManagerHandler_test.ReadResponse"
	fmlogger.Enter(method)

	err := json.NewDecoder(w.Body).Decode(&j)

	if err != nil {
		fmlogger.ExitError(method, "unexpected error", err)
		return err
	}

	fmlogger.Exit(method)
	return nil
}

func TestLoadApiKeyFromFile(t *testing.T) {
	method := "FinanceManagerHandler_test.TestLoadApiKeyFromFile"
	fmlogger.Enter(method)

	//Write File to read from
	pwd, _ := os.Getwd()
	fileName := "TestLoadApiKeyFromFile.keytest"
	content := "test content"

	err := os.WriteFile(pwd+fileName, []byte(content), 0666)

	if err != nil {
		t.Errorf("failed to persist file prior to test")
	}

	fmCp := fmh
	fmCp.StocksApiKeyFileName = fileName
	fmCp.StocksEnabled = false
	fmCp.PolygonApiKey = ""

	//Run Test
	err = fmCp.LoadApiKeyFromFile()
	assert.Nil(t, err)
	assert.True(t, fmCp.StocksEnabled)
	assert.Equal(t, content, fmCp.PolygonApiKey)

	//Run failing test
	fmCp.StocksApiKeyFileName = "someotherfile.keytest"
	err = fmCp.LoadApiKeyFromFile()
	assert.NotNil(t, err)

	//Clean up test file
	err = os.Remove(pwd + fileName)

	if err != nil {
		t.Errorf("failed to clean up test files after test completion")
	}

	fmlogger.Exit(method)
}

func TestUpdateAndPersistAPIKey(t *testing.T) {
	method := "FinanceManagerHandler_test.TestLoadApiKeyFromFile"
	fmlogger.Enter(method)

	pwd, _ := os.Getwd()
	fileName := "TestUpdateAndPersistAPIKey.keytest"
	content := "test content"

	fmCp := fmh
	fmCp.StocksApiKeyFileName = fileName
	fmCp.StocksEnabled = false
	fmCp.PolygonApiKey = ""

	//Run Test
	err := fmCp.UpdateAndPersistAPIKey(content)
	assert.Nil(t, err)
	assert.True(t, fmCp.StocksEnabled)
	assert.Equal(t, content, fmCp.PolygonApiKey)

	//Verify that test was successful
	_, err = os.ReadFile(pwd + fileName)
	assert.Nil(t, err)

	//Clean up the test file
	err = os.Remove(pwd + fileName)

	if err != nil {
		t.Errorf("failed to clean up test files after test completion")
	}

	fmlogger.Exit(method)
}
