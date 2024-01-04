package dbrepo

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"strings"
	"testing"
	"time"
)

func TestGetAllUserCreditCards(t *testing.T) {
	method := "creditcards_test.TestGetAllUserCreditCards"
	fmlogger.Enter(method)

	cc1 := models.CreditCard{ID: 1, UserID: 1, Name: "loan1", Balance: 1000.0, APR: 0.26, MinPayment: 35.00, MinPaymentPercentage: 0.1, CreateDt: time.Now(), LastUpdateDt: time.Now()}
	cc2 := models.CreditCard{ID: 2, UserID: 1, Name: "loan2", Balance: 1000.0, APR: 0.26, MinPayment: 35.00, MinPaymentPercentage: 0.1, CreateDt: time.Now(), LastUpdateDt: time.Now()}
	cc3 := models.CreditCard{ID: 3, UserID: 2, Name: "loan1", Balance: 1000.0, APR: 0.26, MinPayment: 35.00, MinPaymentPercentage: 0.1, CreateDt: time.Now(), LastUpdateDt: time.Now()}

	p.GormDB.Create(&cc1)
	p.GormDB.Create(&cc2)
	p.GormDB.Create(&cc3)

	ccs, err := d.GetAllUserCreditCards(1, "")

	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		t.Errorf(err.Error())
	}

	if len(ccs) != 2 {
		t.Errorf("incorrect length of response items returned")
	}

	ccs, err = d.GetAllUserCreditCards(1, "1")
	if err != nil {
		fmlogger.ExitError(method, err.Error(), err)
		t.Errorf(err.Error())
	}

	if len(ccs) != 1 {
		t.Errorf("incorrect length of response items returned")
	}

	fmlogger.Exit(method)
}

func TestInsertCreditCard(t *testing.T) {
	method := "creditcards_test.TestInsertCreditCard"
	fmlogger.Enter(method)

	p.GormDB.Exec("DELETE FROM credit_cards")
	name := "TestInsertCreditCard"

	cc := models.CreditCard{
		ID:                   1,
		UserID:               1,
		Name:                 name,
		Balance:              1000.0,
		APR:                  26.2,
		MinPayment:           35,
		MinPaymentPercentage: 10,
		CreateDt:             time.Now(),
		LastUpdateDt:         time.Now(),
	}

	id, err := d.InsertCreditCard(cc)

	if err != nil {
		t.Errorf("Unexpected error occured %s", err)
	}

	if id == 0 {
		t.Errorf("ID returned with incorrect value of 0")
	}

	var ccDb models.CreditCard
	p.GormDB.First(&ccDb, id)

	if strings.Compare(ccDb.Name, name) != 0 {
		t.Errorf("Value commited to DB does not match expectations. Expected %s but found %s", name, ccDb.Name)
	}

	fmlogger.Exit(method)
}
