package testingutils

import (
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/models"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const postgres_version = "14.5"
const db_username = "postgres"
const db_password = "postgres"
const db_name = "financemanager_test"

type DockerDBConfig struct {
	Postgres_version string
	DB_username      string
	DB_password      string
	DB_name          string
}

func GetDefaultConfig() DockerDBConfig {
	config := DockerDBConfig{
		Postgres_version: postgres_version,
		DB_username:      db_username,
		DB_password:      db_password,
		DB_name:          db_name,
	}

	return config
}

func InitTables(port string) {
	method := "dockerdb.InitTables"
	fmlogger.Enter(method)

	//Connect to DB
	fmlogger.Info(method, "attemting to connect to DB via GORM")
	config := GetDefaultConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", config.DB_username, config.DB_password, port, config.DB_name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to db")
	}

	fmlogger.Info(method, "successfully connected to DB")

	//Init Tables
	fmlogger.Info(method, "initializing tables")
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Role{})
	db.AutoMigrate(&models.UserRole{})
	db.AutoMigrate(&models.Bill{})
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

	db.Create(&models.User{ID: 1, Username: "admin1", FirstName: "admin", LastName: "admin", Password: "password", Email: "admin@fm.com", CreateDt: time.Now(), LastUpdateDt: time.Now()})
	db.Create(&models.User{ID: 2, Username: "user1", FirstName: "user", LastName: "user", Password: "password", Email: "user@fm.com", CreateDt: time.Now(), LastUpdateDt: time.Now()})

	fmlogger.Exit(method)
}

func seedRoles(db *gorm.DB) {
	method := "dockerdb.seedRoles"
	fmlogger.Enter(method)

	db.Create(&models.Role{ID: 1, Code: "admin", CreateDt: time.Now(), LastUpdateDt: time.Now()})
	db.Create(&models.Role{ID: 2, Code: "user", CreateDt: time.Now(), LastUpdateDt: time.Now()})

	fmlogger.Exit(method)
}

func seedUserRoles(db *gorm.DB) {
	method := "dockerdb.seedUserRoles"
	fmlogger.Enter(method)

	db.Create(&models.UserRole{ID: 1, UserId: 1, RoleId: 1, Code: "admin", CreateDt: time.Now(), LastUpdateDt: time.Now()})

	fmlogger.Exit(method)
}
