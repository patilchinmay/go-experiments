package db

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/sethvargo/go-envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DBName   string `env:"DB_NAME,required"`
}

type Database struct {
	databaseConfig DatabaseConfig
	logger         zerolog.Logger
	DB             *gorm.DB
}

var database *Database

func New(logger zerolog.Logger) *Database {
	if database == nil {
		database = &Database{
			databaseConfig: DatabaseConfig{},
			logger:         logger,
		}

		database.loadCredentials()
		database.connect()
		database.setupConnectionPool()
	}
	return database
}

func (d *Database) loadCredentials() {
	// Load credentials from env vars
	// Uses https://github.com/sethvargo/go-envconfig
	if err := envconfig.Process(context.Background(), &d.databaseConfig); err != nil {
		d.logger.Fatal().Err(err).Msg("Failed to override from env vars")
	}
}

func (d *Database) connect() {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", d.databaseConfig.User, d.databaseConfig.Password, d.databaseConfig.Host, d.databaseConfig.Port, d.databaseConfig.DBName)

	var err error
	d.DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbURL,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		d.logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	d.logger.Debug().Msg("Connected to database")
}

func (d *Database) setupConnectionPool() {
	sqlDB, err := d.DB.DB()
	if err != nil {
		d.logger.Fatal().Err(err).Msg("Failed to get database instance")
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
}
