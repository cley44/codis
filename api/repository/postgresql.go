package repository

import (
	"codis/config"
	"codis/utils"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/samber/lo"

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
		utils.PrintJSONIndent(err.Error())
		panic("Failed to init db")
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
