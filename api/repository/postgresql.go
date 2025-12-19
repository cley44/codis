package repository

import (
	"codis/config"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"github.com/samber/oops"

	// Do not remove: this line imports the driver.
	_ "github.com/lib/pq"
)

type PostgresDatabaseService struct {
	config *config.ConfigService
	Db     *sqlx.DB
}

func NewPostgresDatabaseService(injector do.Injector) (*PostgresDatabaseService, error) {
	config := do.MustInvoke[*config.ConfigService](injector)

	uri := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.Hostname,
		config.Postgres.Port,
		config.Postgres.Database,
		lo.Ternary(config.Postgres.SSL, "require", "disable"),
	)

	db, err := sqlx.Open("postgres", uri)

	if err != nil {
		slog.Error("Failed to init DB", "error", oops.Wrap(err))
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(time.Hour)

	d := PostgresDatabaseService{
		config: config,
		Db:     db,
	}

	return &d, nil
}

func (svc PostgresDatabaseService) Get(dest interface{}, query string, args ...interface{}) error {
	err := svc.Db.Get(dest, query, args...)
	return oops.Wrap(err)
}
