package internal_test

import (
	accounts "assignments/simplebank/accounts/api"
	"assignments/simplebank/accounts/internal"
	"assignments/simplebank/adapters/postgres"
	"os"
	"testing"

	"github.com/kelseyhightower/envconfig"
)

// var testDB *sql.DB
var store accounts.Store

func TestMain(m *testing.M) {
	config := postgres.Config{}
	envconfig.MustProcess("", &config)

	db, testDB := postgres.Connect(config)
	defer testDB.Close()

	store = internal.NewStore(db)
	os.Exit(m.Run())
}
