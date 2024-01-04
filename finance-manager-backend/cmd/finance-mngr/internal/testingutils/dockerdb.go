package testingutils

import (
	"database/sql"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"fmt"
	"log"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/gorm"
)

const postgres_version = "14.5"
const db_username = "postgres"
const db_password = "postgres"
const db_name = "financemanager_test"

var dockerConfig = GetDefaultConfig()

type DockerDBConfig struct {
	Postgres_version string
	DB_username      string
	DB_password      string
	DB_name          string
}

type DockerTestPlatform struct {
	DB       *sql.DB
	GormDB   *gorm.DB
	Resource *dockertest.Resource
	Pool     *dockertest.Pool
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

func Setup(m *testing.M) DockerTestPlatform {
	method := "dockerdb.Setup"
	fmlogger.Enter(method)

	var err error
	var pool *dockertest.Pool
	var resource *dockertest.Resource
	var testDB *sql.DB

	options := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        dockerConfig.Postgres_version,
		Env: []string{
			"POSTGRES_USER=" + dockerConfig.DB_username,
			"POSTGRES_PASSWORD=" + dockerConfig.DB_password,
			"POSTGRES_DB=" + dockerConfig.DB_name,
			"listen_addresses='*'",
		},
	}

	//Create Docker Pool
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not construct docker pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not donnect to Docker: %s", err)
	}

	//Pulls the docker image and creates a container
	resource, err = pool.RunWithOptions(
		&options,
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s timezone=UTC connect_timeout=5", resource.GetPort("5432/tcp"), dockerConfig.DB_username, dockerConfig.DB_password, dockerConfig.DB_name))

		if err != nil {
			return err
		}

		return testDB.Ping()

	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	fmlogger.Info(method, "DB Connection complete")
	gormDb := GetGormDB(resource.GetPort("5432/tcp"))
	plat := DockerTestPlatform{
		DB:       testDB,
		GormDB:   gormDb,
		Resource: resource,
		Pool:     pool,
	}

	fmlogger.Info(method, "Initializing tables and seeding test data")
	InitTables(plat.GormDB)
	fmlogger.Info(method, "Table initialization and test data seeding complete")

	fmlogger.Exit(method)
	return plat
}

func TearDown(p DockerTestPlatform) {
	method := "dockerdb.TearDown"
	fmlogger.Enter(method)

	if err := p.Pool.Purge(p.Resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	fmlogger.Exit(method)
}
