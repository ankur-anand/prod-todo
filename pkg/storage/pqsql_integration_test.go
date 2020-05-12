// +build integration_tests all_tests

package storage_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/ankur-anand/prod-app/pkg/storage"

	"github.com/ankur-anand/prod-app/pkg/storage/auths/testsuite"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

var (
	pgURL *url.URL
	repo  storage.PostgreSQL
)

func TestMain(m *testing.M) {
	code := 0
	pgURL = &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("myuser", "mypass"),
		Path:   "userdatabase",
	}
	q := pgURL.Query()
	q.Add("sslmode", "disable")
	pgURL.RawQuery = q.Encode()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("cloud not connect with docker API %v", err)
	}

	password, _ := pgURL.User.Password()
	// Pull an image, create a container based on it and set all necessary parameters
	runOpts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=" + pgURL.User.Username(),
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + pgURL.Path,
		},
	}

	resource, err := pool.RunWithOptions(&runOpts)
	if err != nil {
		log.Fatalf("cloud not start postgres container %v", err)
	}

	// container IP address from host is the DB IP
	pgURL.Host = resource.Container.NetworkSettings.IPAddress

	// network on Mac
	if runtime.GOOS == "darwin" {
		pgURL.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}

	logWaiter, err := pool.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    resource.Container.ID,
		OutputStream: log.Writer(),
		ErrorStream:  log.Writer(),
		Stderr:       true,
		Stdout:       true,
		Stream:       true,
	})
	if err != nil {
		log.Fatalf("could not connect to postgres container log output %v", err)
	}
	defer func() {
		err = logWaiter.Close()
		if err != nil {
			log.Println("ERROR: could not close container log")
		}
		err = logWaiter.Wait()
		if err != nil {
			log.Println("error: could not wait for container log to close")
		}
	}()

	pool.MaxWait = 10 * time.Second
	var conn *pgx.Conn
	err = pool.Retry(func() error {
		// connect to the db
		conn, err = pgx.Connect(context.Background(), pgURL.String())
		if err != nil {
			return err
		}
		return conn.Ping(context.Background())
	})

	if err != nil {
		log.Fatal("could not connect to postgres server")
	}

	// migrate
	_, filename, _, _ := runtime.Caller(0)
	fDirName := filepath.Dir(filename)
	migrationDir := filepath.Join(fDirName, "migrations")
	sourceURL := fmt.Sprintf("file://%s", migrationDir)
	db, err := sql.Open("postgres", pgURL.String())
	if err != nil {
		log.Fatalf("error opening postgres connection %v", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("error getting postgres driver instance %v", err)
	}
	// https://github.com/golang-migrate/migrate/issues/226
	mig, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres", driver)
	if err != nil {
		log.Fatalf("error getting migrate down instance %v", err)
	}

	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	}

	repo, err = storage.NewPostgreSQL(pgURL.String())
	if err != nil {
		log.Fatalf("error init repo connection %v", err)
	}
	// start other test
	code = m.Run()
	err = conn.Close(context.Background())
	if err != nil {
		log.Println("error closing db connection")
	}
	err = pool.Purge(resource)
	if err != nil {
		log.Println("error could not purge resource")
	}
	os.Exit(code)
}

func TestFindAndStore(t *testing.T) {
	t.Parallel()
	suiteBase := &testsuite.SuiteBase{}
	suiteBase.SetRepo(repo.AuthSQL())
	suiteBase.TestFindAndStore(t)
}

func TestFindByEmailAndStore(t *testing.T) {
	t.Parallel()
	suiteBase := &testsuite.SuiteBase{}
	suiteBase.SetRepo(repo.AuthSQL())
	suiteBase.TestFindByEmailAndStore(t)
}

func TestDuplicateEmailStorePqSQL(t *testing.T) {
	t.Parallel()
	suiteBase := &testsuite.SuiteBase{}
	suiteBase.SetRepo(repo.AuthSQL())
	suiteBase.TestDuplicateEmailStorePqSQL(t)
}
