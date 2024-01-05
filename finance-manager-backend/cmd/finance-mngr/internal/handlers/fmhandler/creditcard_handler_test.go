package fmhandler

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"finance-manager-backend/cmd/finance-mngr/internal/testingutils"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAllUserCreditCards_200(t *testing.T) {
	method := "creditcard_handler_test.TestGetAllUserCreditCards_200"
	fmlogger.Enter(method)

	setup()
	token := testingutils.GetAdminJWT(t)

	writer := MakeRequest(http.MethodGet, "/users/1/credit-cards", nil, true, token)

	assert.Equal(t, http.StatusOK, writer.Code)

	tearDown()
	fmlogger.Exit(method)
}

func TestGetAllUserCreditCards_403(t *testing.T) {
	method := "creditcard_handler_test.TestGetAllUserCreditCards_403"
	fmlogger.Enter(method)

	token := testingutils.GetUserJWT(t)

	writer := MakeRequest(http.MethodGet, "/users/1/credit-cards", nil, true, token)

	assert.Equal(t, http.StatusForbidden, writer.Code)

	fmlogger.Exit(method)
}

func TestSaveCreditCard_200(t *testing.T) {
	method := "creditcard_handler_test.TestSaveCreditCard"
	fmlogger.Enter(method)

	var ccDb models.CreditCard
	cc := models.CreditCard{
		UserID:               1,
		Name:                 "TestSaveCreditCard",
		Balance:              1000.0,
		Limit:                20000.0,
		APR:                  26.2,
		MinPayment:           35.00,
		MinPaymentPercentage: 10,
	}

	token := testingutils.GetAdminJWT(t)

	writer := MakeRequest(http.MethodPost, "/users/1/credit-cards", cc, true, token)

	assert.Equal(t, http.StatusAccepted, writer.Code)

	//Validate object was saved into the database
	p.GormDB.Where("name = ?", cc.Name).First(&ccDb)
	assert.Equal(t, cc.Name, ccDb.Name)
	assert.Equal(t, cc.Balance, ccDb.Balance)
	assert.NotEqual(t, cc.CreateDt, ccDb.CreateDt)

	//Clear the DB
	tearDown()

	fmlogger.Exit(method)
}

func TestSaveCreditCard_403(t *testing.T) {
	method := "creditcard_handler_test.TestSaveCreditCard_403"
	fmlogger.Enter(method)

	cc := models.CreditCard{
		UserID:               1,
		Name:                 "TestSaveCreditCard",
		Balance:              1000.0,
		Limit:                20000.0,
		APR:                  26.2,
		MinPayment:           35.00,
		MinPaymentPercentage: 10,
	}

	token := testingutils.GetUserJWT(t)

	writer := MakeRequest(http.MethodPost, "/users/1/credit-cards", cc, true, token)

	assert.Equal(t, http.StatusForbidden, writer.Code)

	fmlogger.Exit(method)
}

func TestSaveCreditCard_400(t *testing.T) {
	method := "creditcard_handler_test.TestSaveCreditCard_400"
	fmlogger.Enter(method)
	cc := models.CreditCard{
		UserID:               0,
		Name:                 "",
		Balance:              1000.0,
		Limit:                20000.0,
		APR:                  26.2,
		MinPayment:           35.00,
		MinPaymentPercentage: 10,
	}

	token := testingutils.GetAdminJWT(t)

	//Empty Object
	writer := MakeRequest(http.MethodPost, "/users/1/credit-cards", nil, true, token)
	assert.Equal(t, http.StatusBadRequest, writer.Code)

	//Malformed Object
	writer = MakeRequest(http.MethodPost, "/users/1/credit-cards", "{Bad", true, token)
	assert.Equal(t, http.StatusBadRequest, writer.Code)

	//Empty Name
	writer = MakeRequest(http.MethodPost, "/users/1/credit-cards", cc, true, token)
	assert.Equal(t, http.StatusBadRequest, writer.Code)

	fmlogger.Exit(method)
}

func TestGetCreditCardById_200(t *testing.T) {
	method := "creditcard_handler_test.TestGetCreditCardById_200"
	fmlogger.Enter(method)

	setup()
	token := testingutils.GetUserJWT(t)

	writer := MakeRequest(http.MethodGet, "/users/2/credit-cards/3", nil, true, token)
	assert.Equal(t, http.StatusOK, writer.Code)

	tearDown()
	fmlogger.Exit(method)
}

func TestGetCreditCardById_403(t *testing.T) {
	method := "creditcard_handler_test.TestGetCreditCardById_403"
	fmlogger.Enter(method)

	setup()
	token := testingutils.GetUserJWT(t)

	//Credit Card belongs to other user
	writer := MakeRequest(http.MethodGet, "/users/2/credit-cards/1", nil, true, token)
	assert.Equal(t, http.StatusForbidden, writer.Code)

	//userId is for other user
	writer = MakeRequest(http.MethodGet, "/users/1/credit-cards/1", nil, true, token)
	assert.Equal(t, http.StatusForbidden, writer.Code)

	tearDown()
	fmlogger.Exit(method)
}

func TestGetCreditCardById_400(t *testing.T) {
	method := "creditcard_handler_test.TestGetCreditCardById_403"
	fmlogger.Enter(method)

	setup()
	token := testingutils.GetUserJWT(t)

	//Invalid ID
	writer := MakeRequest(http.MethodGet, "/users/2/credit-cards/a", nil, true, token)
	assert.Equal(t, http.StatusBadRequest, writer.Code)

	tearDown()
	fmlogger.Exit(method)
}

func TestGetCreditCardById_404(t *testing.T) {
	method := "creditcard_handler_test.TestGetCreditCardById_422"
	fmlogger.Enter(method)

	setup()
	token := testingutils.GetUserJWT(t)

	//Invalid ID
	writer := MakeRequest(http.MethodGet, "/users/2/credit-cards/4", nil, true, token)
	assert.Equal(t, http.StatusNotFound, writer.Code)

	tearDown()
	fmlogger.Exit(method)
}

func setup() {
	method := "creditcard_handler_test.setup"
	fmlogger.Enter(method)

	cc1 := models.CreditCard{ID: 1, UserID: 1, Name: "cc1", Balance: 1000.0, APR: 0.26, MinPayment: 35.00, MinPaymentPercentage: 0.1, CreateDt: time.Now(), LastUpdateDt: time.Now()}
	cc2 := models.CreditCard{ID: 2, UserID: 1, Name: "cc2", Balance: 1000.0, APR: 0.26, MinPayment: 35.00, MinPaymentPercentage: 0.1, CreateDt: time.Now(), LastUpdateDt: time.Now()}
	cc3 := models.CreditCard{ID: 3, UserID: 2, Name: "cc3", Balance: 1000.0, APR: 0.26, MinPayment: 35.00, MinPaymentPercentage: 0.1, CreateDt: time.Now(), LastUpdateDt: time.Now()}

	p.GormDB.Create(&cc1)
	p.GormDB.Create(&cc2)
	p.GormDB.Create(&cc3)
	fmlogger.Info(method, "Inserted credit card test data")
	fmlogger.Exit(method)
}

func tearDown() {
	method := "creditcard_handler_test.tearDown"
	fmlogger.Enter(method)

	//Clean Data
	p.GormDB.Exec("DELETE FROM credit_cards")

	fmlogger.Info(method, "Cleaned all credit card test data")
	fmlogger.Exit(method)
}
