package repository

import (
	"codis/config"
	"time"

	"github.com/samber/do/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabaseService struct {
	config *config.ConfigService
	db     *gorm.DB
}

func NewPostgresDatabaseService(injector do.Injector) (*PostgresDatabaseService, error) {
	config := do.MustInvoke[*config.ConfigService](injector)

	// https://github.com/go-gorm/postgres
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  config.Postgres.URI,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		panic("Failed to init db")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("sql db failed")
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	d := PostgresDatabaseService{
		config: config,
		db:     db,
	}

	return &d, nil
}
