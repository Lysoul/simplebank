package internal_test

import (
	accounts "assignments/simplebank/accounts/api"
	"assignments/simplebank/accounts/internal"
	"assignments/simplebank/adapters/monitoring"
	"assignments/simplebank/adapters/postgres"
	"os"
	"testing"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

// var testDB *sql.DB
var store accounts.Store

func TestMain(m *testing.M) {
	config := postgres.Config{}
	envconfig.MustProcess("", &config)

	monitoring.Logger().Info("Starting tests with database ", zap.String("url", config.URL))
	db, testDB := postgres.Connect(config)
	defer testDB.Close()

	store = internal.NewStore(db)
	os.Exit(m.Run())
}
