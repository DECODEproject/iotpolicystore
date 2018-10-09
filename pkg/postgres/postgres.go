package postgres

import (
	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/thingful/iotpolicystore/pkg/config"
)

// DB is an exported type that exposes methods for reading/writing to our
// postgres database.
type DB struct {
	DB *sqlx.DB

	connStr string
	verbose bool
	logger  kitlog.Logger
}

// Open is a helper function that connects to the DB and returns an instantiated
// sqlx.DB instance or an error.
func Open(connStr string) (*sqlx.DB, error) {
	return sqlx.Open("postgres", connStr)
}

// NewDB is a constructor for an instance of our DB type. Sets up the logger,
// and other configuration details, but does not connect to the DB until Start()
// is called.
func NewDB(config *config.Config) *DB {
	logger := kitlog.With(config.Logger, "module", "postgres")

	if config.Verbose {
		logger.Log("msg", "creating postgres DB connection")
	}

	return &DB{
		connStr: config.ConnStr,
		verbose: config.Verbose,
		logger:  logger,
	}
}

// Start tells the DB component to actually connect to the database and be ready
// to start work.
func (d *DB) Start() error {
	d.logger.Log("msg", "starting postgres DB connection")

	db, err := Open(d.connStr)
	if err != nil {
		return errors.Wrap(err, "failed to open DB connection")
	}

	d.DB = db

	return nil
}

// Stop tells the DB component to close the DB connection pool.
func (d *DB) Stop() error {
	d.logger.Log("msg", "stopping postgres DB connection")

	if d.DB != nil {
		return d.DB.Close()
	}

	return nil
}
