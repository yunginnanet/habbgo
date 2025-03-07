package test

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

const dbImageName = "mariadb"

var db *sql.DB

func TestDatabaseSetup(t *testing.T) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	// pulls an image, creates a container based on it and runs it
	t.Logf("Building database container with image: %s...", dbImageName)
	resource, err := pool.Run(dbImageName, "latest", []string{"MYSQL_ROOT_PASSWORD=secret"})
	require.NoError(t, err)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	t.Log("Connectiong to the database...")
	require.NoError(t, pool.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}))

}
