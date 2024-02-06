package test

import (
	"finance-manager-backend/internal/finance-mngr/models"
	"fmt"
	"time"

	"github.com/jon-kamis/klogger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TestingAdmin = models.User{ID: 1, Username: "admin1", FirstName: "admin", LastName: "admin", Password: "password", Email: "admin@fm.com", CreateDt: time.Now(), LastUpdateDt: time.Now()}
var TestingUser = models.User{ID: 2, Username: "user1", FirstName: "user", LastName: "user", Password: "password", Email: "user@fm.com", CreateDt: time.Now(), LastUpdateDt: time.Now()}
var AdminRole = models.Role{ID: 1, Code: "admin", CreateDt: time.Now(), LastUpdateDt: time.Now()}
var UserRole = models.Role{ID: 2, Code: "user", CreateDt: time.Now(), LastUpdateDt: time.Now()}

func GetGormDB(port string) *gorm.DB {
	method := "dockerdb.GetGormDB"
	klogger.Enter(method)

	klogger.Info(method, "attemting to connect to DB via GORM")
	config := GetDefaultConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", config.DB_username, config.DB_password, port, config.DB_name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to db")
	}

	klogger.Info(method, "successfully connected to DB")

	klogger.Exit(method)
	return db
}

func InitTables(db *gorm.DB) {
	method := "dockerdb.InitTables"
	klogger.Enter(method)

	//Init Tables
	klogger.Info(method, "initializing tables")
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Role{})
	db.AutoMigrate(&models.UserRole{})
	db.AutoMigrate(&models.Bill{})
	db.AutoMigrate(&models.CreditCard{})
	db.AutoMigrate(&models.Stock{})
	db.AutoMigrate(&models.UserStock{})
	db.AutoMigrate(&models.StockData{})
	klogger.Info(method, "tables initialized")

	//Seed Data
	klogger.Info(method, "seeding data")
	seedUsers(db)
	seedRoles(db)
	seedUserRoles(db)
	klogger.Info(method, "data seeded")

	klogger.Info(method, "db initialization successful")
	klogger.Exit(method)
}

func seedUsers(db *gorm.DB) {
	method := "dockerdb.seedUsers"
	klogger.Enter(method)

	db.Create(&TestingAdmin)
	db.Create(&TestingUser)

	klogger.Exit(method)
}

func seedRoles(db *gorm.DB) {
	method := "dockerdb.seedRoles"
	klogger.Enter(method)

	db.Create(&AdminRole)
	db.Create(&UserRole)

	klogger.Exit(method)
}

func seedUserRoles(db *gorm.DB) {
	method := "dockerdb.seedUserRoles"
	klogger.Enter(method)

	db.Create(&models.UserRole{ID: 1, UserId: TestingAdmin.ID, RoleId: 1, Code: "admin", CreateDt: time.Now(), LastUpdateDt: time.Now()})
	db.Create(&models.UserRole{ID: 2, UserId: TestingAdmin.ID, RoleId: 2, Code: "user", CreateDt: time.Now(), LastUpdateDt: time.Now()})
	db.Create(&models.UserRole{ID: 3, UserId: TestingUser.ID, RoleId: 2, Code: "user", CreateDt: time.Now(), LastUpdateDt: time.Now()})

	klogger.Exit(method)
}
