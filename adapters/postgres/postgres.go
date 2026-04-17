package postgres

import (
	accountApi "assignments/simplebank/accounts/api"
	"assignments/simplebank/adapters/monitoring"
	"context"
	"database/sql"
	"regexp"
	"runtime"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	URL     string `envconfig:"POSTGRES_URL" required:"true"`
	Migrate bool   `envconfig:"POSTGRES_MIGRATE"`
	Debug   bool   `envconfig:"POSTGRES_DEBUG" default:"true"`
	// default to 4 * runtime.NumCPU
	MaxOpenConns int `envconfig:"POSTGRES_MAX_OPEN_CONNS"`
	// default to 4 * runtime.NumCPU
	MaxIdleConns int `envconfig:"POSTGRES_MAX_IDLE_CONNS"`
	// To make your app more resilient to errors during migrations, you can tweak Bun to discard unknown columns in production
	DiscardUnknownColumns bool          `envconfig:"POSTGRES_BUN_DISCARD_UNKNOWN_COLUMNS"`
	SlowQueriesDuration   time.Duration `envconfig:"POSTGRES_SLOW_QUERIES_DURATION" default:"200ms"`
}

func Connect(config Config) (*gorm.DB, *sql.DB) {

	log := monitoring.Logger()
	log.Info("Connecting to database " +
		regexp.MustCompile(`postgres://.+@`).
			ReplaceAllString(config.URL, "postgres://xxxx:xxxx@"))

	if config.MaxOpenConns == 0 && config.MaxIdleConns == 0 {
		maxOpenConns := 4 * runtime.GOMAXPROCS(0)
		config.MaxOpenConns = maxOpenConns
		config.MaxIdleConns = maxOpenConns
	}

	gormConfig := &gorm.Config{}
	sqldb, err := gorm.Open(postgres.Open(config.URL), gormConfig)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	db, err := sqldb.DB()
	if err != nil {
		log.Fatal("failed to get database instance", zap.Error(err))
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	if config.Migrate {
		//todo: improve with better approach for model registration and migration
		//we can also split this as the cli e.g. `simplebank db migrate` and `simplebank db rollback`
		err := sqldb.AutoMigrate(&accountApi.Account{}, &accountApi.Transfer{}, &accountApi.Entry{})
		if err != nil {
			log.Fatal("failed to migrate", zap.Error(err))
		}
		log.Info("migration completed")
	}

	return sqldb, db

}

func HealthFunc(db *gorm.DB) func(context.Context) error {
	return func(ctx context.Context) error {
		return db.WithContext(ctx).Raw("SELECT 1").Error
	}
}
