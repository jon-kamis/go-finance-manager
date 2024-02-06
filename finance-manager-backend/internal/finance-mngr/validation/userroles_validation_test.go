package validation

import (
	"finance-manager-backend/internal/finance-mngr/models"
	"finance-manager-backend/test"
	"testing"
	"time"

	"github.com/jon-kamis/klogger"
)

func TestUserRoleBelongsToUser(t *testing.T) {
	method := "userroles_validation_test.TestUserRoleBelongsToUser"
	klogger.Enter(method)

	userId := test.TestingAdmin.ID

	userRole := models.UserRole{
		ID:           1,
		UserId:       userId,
		RoleId:       test.AdminRole.ID,
		Code:         test.AdminRole.Code,
		CreateDt:     time.Now(),
		LastUpdateDt: time.Now(),
	}

	err := fmv.UserRoleBelongsToUser(models.UserRole{}, 1)

	if err == nil {
		t.Errorf("expected error for empty userRole object but none was thrown")
	}

	err = fmv.UserRoleBelongsToUser(userRole, 0)

	if err == nil {
		t.Errorf("expected error for invalid userId but none was thrown")
	}

	err = fmv.UserRoleBelongsToUser(userRole, 23)

	if err == nil {
		t.Errorf("expected error for userrole not belonging to user but none was thrown")
	}

	err = fmv.UserRoleBelongsToUser(userRole, userId)

	if err != nil {
		t.Errorf("unexpected error was thrown for valid case")
	}

	klogger.Exit(method)

}

func TestUserRoleExistsAndBelongsToUser(t *testing.T) {
	method := "userroles_validation_test.TestUserRoleExistsAndBelongsToUser"
	klogger.Enter(method)

	userId := 23

	userRole := models.UserRole{
		ID:           23,
		UserId:       23,
		RoleId:       test.AdminRole.ID,
		Code:         test.AdminRole.Code,
		CreateDt:     time.Now(),
		LastUpdateDt: time.Now(),
	}

	//Save into the database
	p.GormDB.Create(&userRole)

	err := fmv.UserRoleExistsAndBelongsToUser(userRole.RoleId, userId)

	if err != nil {
		t.Errorf("unexpected error was thrown for valid case")
	}

	err = fmv.UserRoleExistsAndBelongsToUser(24, userId)

	if err == nil {
		t.Errorf("expected an error for user role does not exist but none was thrown")
	}

	err = fmv.UserRoleExistsAndBelongsToUser(userRole.ID, 24)

	if err == nil {
		t.Errorf("expected an error for user role does not belong to user but none was thrown")
	}

	//Cleanup DB column
	p.GormDB.Delete(&userRole)

	klogger.Exit(method)

}
