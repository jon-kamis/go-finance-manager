package validation

import (
	"database/sql"
	"finance-manager-backend/cmd/finance-mngr/internal/fmlogger"
	"finance-manager-backend/cmd/finance-mngr/internal/repository/dbrepo"
	"finance-manager-backend/cmd/finance-mngr/internal/testingutils"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var fmv FinanceManagerValidator

var dockerConfig = testingutils.GetDefaultConfig()

func TestMain(m *testing.M) {
	method := "validation_test.go"
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

	fmv = FinanceManagerValidator{DB: &dbrepo.PostgresDBRepo{DB: testDB}}
	fmlogger.Info(method, "DB Connection complete")

	fmlogger.Info(method, "Initializing tables and seeding test data")
	testingutils.InitTables(resource.GetPort("5432/tcp"))
	fmlogger.Info(method, "Table initialization and test data seeding complete")

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	fmlogger.Exit(method)
	os.Exit(code)
}
