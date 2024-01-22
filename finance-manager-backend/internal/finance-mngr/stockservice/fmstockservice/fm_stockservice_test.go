package fmstockservice

import (
	"finance-manager-backend/internal/finance-mngr/fmlogger"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	fss := FmStockService{
		StocksApiKeyFileName : fileName,
		StocksEnabled: false,
	}

	//Run Test
	err = fss.LoadApiKeyFromFile()
	assert.Nil(t, err)
	assert.True(t, fss.StocksEnabled)
	assert.Equal(t, content, fss.PolygonApiKey)

	//Run failing test
	fss.StocksApiKeyFileName = "someotherfile.keytest"
	err = fss.LoadApiKeyFromFile()
	assert.NotNil(t, err)

	//Clean up test file
	err = os.Remove(pwd + fileName)

	if err != nil {
		t.Errorf("failed to clean up test files after test completion")
	}

	fmlogger.Exit(method)
}

func TestUpdateAndPersistAPIKey(t *testing.T) {
	method := "fm_stockservice.TestLoadApiKeyFromFile"
	fmlogger.Enter(method)

	pwd, _ := os.Getwd()
	fileName := "TestUpdateAndPersistAPIKey.keytest"
	content := "test content"

	fss := FmStockService{
		StocksApiKeyFileName : fileName,
		StocksEnabled: false,
	}

	//Run Test
	err := fss.UpdateAndPersistAPIKey(content)
	assert.Nil(t, err)
	assert.True(t, fss.StocksEnabled)
	assert.Equal(t, content, fss.PolygonApiKey)

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
