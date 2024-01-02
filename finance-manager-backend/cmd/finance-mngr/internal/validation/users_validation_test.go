package validation

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"testing"
)

func TestIsValidToViewOtherUserData(t *testing.T) {
	method := "users_validation_test.TestIsValidToViewOtherUserData"
	fmlogger.Enter(method)

	isValid, err := fmv.IsValidToViewOtherUserData(1)

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching user roles\n", method)
		fmlogger.Exit(method)
		t.Errorf("unexpected error when fetching user roles: %s\n", err)
	}

	if !isValid {
		fmlogger.Info(method, "unexpected result returned from test")
		fmlogger.Exit(method)
		t.Errorf("unexpected result returned from test")
	}

	isValid, err = fmv.IsValidToViewOtherUserData(2)

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching user roles\n", method)
		fmlogger.Exit(method)
		t.Errorf("unexpected error when fetching user roles: %s\n", err)
	}

	if isValid {
		fmlogger.Info(method, "unexpected result returned from test")
		fmlogger.Exit(method)
		t.Errorf("unexpected result returned from test")
	}

	fmlogger.Exit(method)
}

func TestIsValidToEnterNewUser(t *testing.T) {
	method := "users_validation_test.TestIsValidToEnterNewUser"
	fmlogger.Enter(method)

	var tu models.User

	u := models.User{
		Username:  "TestUsr1",
		FirstName: "Test",
		LastName:  "User",
		Email:     "testusr1@fm.com",
		Password:  "testPswd",
	}

	//Assert Valid case passes
	err := fmv.IsValidToEnterNewUser(u)

	if err != nil {
		errMsg := "unexpected error occured when validating valid user case"
		fmlogger.ExitError(method, errMsg, err)
		t.Errorf(errMsg)
	}

	//Username is required
	tu = u
	tu.Username = ""
	err = fmv.IsValidToEnterNewUser(tu)

	if err == nil {
		errMsg := "expected error for invalid username but none was thrown"
		fmlogger.ExitError(method, errMsg, err)
		t.Errorf(errMsg)
	}

	//Email is required
	tu = u
	tu.Email = ""
	err = fmv.IsValidToEnterNewUser(tu)

	if err == nil {
		errMsg := "expected error for invalid email but none was thrown"
		fmlogger.ExitError(method, errMsg, err)
		t.Errorf(errMsg)
	}

	//Password is required
	tu = u
	tu.Password = ""
	err = fmv.IsValidToEnterNewUser(tu)

	if err == nil {
		errMsg := "expected error for invalid password but none was thrown"
		fmlogger.ExitError(method, errMsg, err)
		t.Errorf(errMsg)
	}

	//First Name is required
	tu = u
	tu.FirstName = ""
	err = fmv.IsValidToEnterNewUser(tu)

	if err == nil {
		errMsg := "expected error for invalid first name but none was thrown"
		fmlogger.ExitError(method, errMsg, err)
		t.Errorf(errMsg)
	}

	//Last Name is required
	tu = u
	tu.LastName = ""
	err = fmv.IsValidToEnterNewUser(tu)

	if err == nil {
		errMsg := "expected error for invalid last name but none was thrown"
		fmlogger.ExitError(method, errMsg, err)
		t.Errorf(errMsg)
	}

	//Username Already Exists
	tu = u
	tu.Username = "admin1"
	err = fmv.IsValidToEnterNewUser(tu)

	if err == nil {
		errMsg := "expected error for username taken but none was thrown"
		fmlogger.ExitError(method, errMsg, err)
		t.Errorf(errMsg)
	}

	//Email Already Exists
	tu = u
	tu.Email = "admin@fm.com"
	err = fmv.IsValidToEnterNewUser(tu)

	if err == nil {
		errMsg := "expected error for email taken but none was thrown"
		fmlogger.ExitError(method, errMsg, err)
		t.Errorf(errMsg)
	}

	fmlogger.Exit(method)
}
