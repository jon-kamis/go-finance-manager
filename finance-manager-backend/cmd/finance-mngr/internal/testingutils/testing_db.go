package testingutils

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TestingAdmin = models.User{ID: 1, Username: "admin1", FirstName: "admin", LastName: "admin", Password: "password", Email: "admin@fm.com", CreateDt: time.Now(), LastUpdateDt: time.Now()}
var TestingUser = models.User{ID: 2, Username: "user1", FirstName: "user", LastName: "user", Password: "password", Email: "user@fm.com", CreateDt: time.Now(), LastUpdateDt: time.Now()}
var AdminRole = models.Role{ID: 1, Code: "admin", CreateDt: time.Now(), LastUpdateDt: time.Now()}
var UserRole = models.Role{ID: 2, Code: "user", CreateDt: time.Now(), LastUpdateDt: time.Now()}

func GetGormDB(port string) *gorm.DB {
	method := "dockerdb.GetGormDB"
	fmlogger.Enter(method)

	fmlogger.Info(method, "attemting to connect to DB via GORM")
	config := GetDefaultConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", config.DB_username, config.DB_password, port, config.DB_name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to db")
	}

	fmlogger.Info(method, "successfully connected to DB")

	fmlogger.Exit(method)
	return db
}

func InitTables(db *gorm.DB) {
	method := "dockerdb.InitTables"
	fmlogger.Enter(method)

	//Init Tables
	fmlogger.Info(method, "initializing tables")
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Role{})
	db.AutoMigrate(&models.UserRole{})
	db.AutoMigrate(&models.Bill{})
	db.AutoMigrate(&models.CreditCard{})
	db.AutoMigrate(&models.Stock{})
	db.AutoMigrate(&models.UserStock{})
	db.AutoMigrate(&models.StockData{})
	fmlogger.Info(method, "tables initialized")

	//Seed Data
	fmlogger.Info(method, "seeding data")
	seedUsers(db)
	seedRoles(db)
	seedUserRoles(db)
	fmlogger.Info(method, "data seeded")

	fmlogger.Info(method, "db initialization successful")
	fmlogger.Exit(method)
}

func seedUsers(db *gorm.DB) {
	method := "dockerdb.seedUsers"
	fmlogger.Enter(method)

	db.Create(&TestingAdmin)
	db.Create(&TestingUser)

	fmlogger.Exit(method)
}

func seedRoles(db *gorm.DB) {
	method := "dockerdb.seedRoles"
	fmlogger.Enter(method)

	db.Create(&AdminRole)
	db.Create(&UserRole)

	fmlogger.Exit(method)
}

func seedUserRoles(db *gorm.DB) {
	method := "dockerdb.seedUserRoles"
	fmlogger.Enter(method)

	db.Create(&models.UserRole{ID: 1, UserId: TestingAdmin.ID, RoleId: 1, Code: "admin", CreateDt: time.Now(), LastUpdateDt: time.Now()})
	db.Create(&models.UserRole{ID: 2, UserId: TestingAdmin.ID, RoleId: 2, Code: "user", CreateDt: time.Now(), LastUpdateDt: time.Now()})
	db.Create(&models.UserRole{ID: 3, UserId: TestingUser.ID, RoleId: 2, Code: "user", CreateDt: time.Now(), LastUpdateDt: time.Now()})

	fmlogger.Exit(method)
}
