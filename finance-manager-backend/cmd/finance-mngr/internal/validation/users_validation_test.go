package validation

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"fmt"
	"testing"
)

func TestCheckIfUserHasRole(t *testing.T) {
	method := "users_validation_test.TestCheckIfUserHasRole"
	fmlogger.Enter(method)

	_, err := fmv.IsValidToViewOtherUserData(1)

	if err != nil {
		fmt.Printf("[%s] unexpected error occured when fetching user roles\n", method)
		fmlogger.Exit(method)
		t.Errorf("unexpected error when fetching user roles: %s\n", err)
	}

	fmlogger.Exit(method)
}
