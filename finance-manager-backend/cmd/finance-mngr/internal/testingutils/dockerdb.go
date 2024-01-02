package testingutils

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
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

type DockerDBTest struct {
	DB       *sql.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (d *DockerDBTest) InitDb() {
	var err error
	var pool *dockertest.Pool
	var resource *dockertest.Resource
	var testDB *sql.DB

	options := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        postgres_version,
		Env: []string{
			"POSTGRES_USER=" + db_username,
			"POSTGRES_PASSWORD=" + db_password,
			"POSTGRES_DB=" + db_name,
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
		testDB, err = sql.Open("pgx", fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s timezone=UTC connect_timeout=5", resource.GetPort("5432/tcp"), db_username, db_password, db_name))

		if err != nil {
			return err
		}

		return testDB.Ping()

	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	d.pool = pool
	d.resource = resource
	d.DB = testDB
}

func (d *DockerDBTest) Destroy() {
	if err := d.pool.Purge(d.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
