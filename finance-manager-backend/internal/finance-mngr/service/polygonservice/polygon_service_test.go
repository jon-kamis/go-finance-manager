package polygonservice

import (
	"finance-manager-backend/internal/finance-mngr/jsonutils"
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/test"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/jon-kamis/klogger"
	"github.com/stretchr/testify/assert"
)

var ps PolygonService

func TestMain(m *testing.M) {
	method := "validation_test.TestMain"
	klogger.Enter(method)

	mockUrl := "http://localhost"
	mockPort := 8081

	mockApi := test.MockPolygonApi{
		BaseUrl: mockUrl,
		Port:    mockPort,
		Handler: test.MockPolygonHandler{
			JSONUtil: &jsonutils.JSONUtil{},
		},
	}

	ps = PolygonService{
		BaseApi:       "http://localhost:8081",
		StocksEnabled: true,
		PolygonApiKey: "test",
	}

	//Start up a mock server in the background
	go startMockApiServer(mockPort, mockApi)

	//Execute Code
	code := m.Run()

	klogger.Exit(method)
	os.Exit(code)

}

func startMockApiServer(p int, m test.MockPolygonApi) {
	method := "polygon_service_test.startMockApiServer"
	klogger.Enter(method)

	//start a web server
	err := http.ListenAndServe(fmt.Sprintf(":%d", p), m.Routes())
	if err != nil {
		log.Fatal(err)
	}

}

func TestLoadApiKeyFromFile(t *testing.T) {
	method := "polygon_service_test.TestLoadApiKeyFromFile"
	klogger.Enter(method)

	//Write File to read from
	pwd, _ := os.Getwd()
	fileName := "TestLoadApiKeyFromFile.keytest"
	content := "test content"

	err := os.WriteFile(pwd+fileName, []byte(content), 0666)

	if err != nil {
		t.Errorf("failed to persist file prior to test")
	}

	fms := PolygonService{
		StocksApiKeyFileName: fileName,
		StocksEnabled:        false,
	}

	//Run Test
	err = fms.LoadApiKeyFromFile()
	assert.Nil(t, err)
	assert.True(t, fms.StocksEnabled)
	assert.Equal(t, content, fms.PolygonApiKey)

	//Run failing test
	fms.StocksApiKeyFileName = "someotherfile.keytest"
	err = fms.LoadApiKeyFromFile()
	assert.NotNil(t, err)

	//Clean up test file
	err = os.Remove(pwd + fileName)

	if err != nil {
		t.Errorf("failed to clean up test files after test completion")
	}

	klogger.Exit(method)
}

func TestUpdateAndPersistAPIKey(t *testing.T) {
	method := "fm_stockservice.TestLoadApiKeyFromFile"
	klogger.Enter(method)

	pwd, _ := os.Getwd()
	fileName := "TestUpdateAndPersistAPIKey.keytest"
	content := "test content"

	fms := PolygonService{
		StocksApiKeyFileName: fileName,
		StocksEnabled:        false,
	}

	//Run Test
	err := fms.UpdateAndPersistAPIKey(content)
	assert.Nil(t, err)
	assert.True(t, fms.StocksEnabled)
	assert.Equal(t, content, fms.PolygonApiKey)

	//Verify that test was successful
	_, err = os.ReadFile(pwd + fileName)
	assert.Nil(t, err)

	//Clean up the test file
	err = os.Remove(pwd + fileName)

	if err != nil {
		t.Errorf("failed to clean up test files after test completion")
	}

	klogger.Exit(method)
}

func TestFetchStockWithTicker(t *testing.T) {
	method := "polygon_service_test.TestFetchStockWithTicker"
	klogger.Enter(method)

	ticker := "AAPL"
	var s models.Stock
	var err error

	s, err = ps.FetchStockWithTicker(ticker)

	assert.Nil(t, err)
	assert.Equal(t, ticker, s.Ticker)
	assert.NotEqual(t, 0.0, s.High)
	assert.NotEqual(t, 0.0, s.Low)
	assert.NotEqual(t, 0.0, s.Open)
	assert.NotEqual(t, 0.0, s.Close)
	assert.False(t, s.Date.IsZero())

	klogger.Exit(method)
}
