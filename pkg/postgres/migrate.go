package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/DECODEproject/iotpolicystore/pkg/migrations"
	kitlog "github.com/go-kit/kit/log"
	"github.com/golang-migrate/migrate"
	psm "github.com/golang-migrate/migrate/database/postgres"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/serenize/snaker"
)

// MigrateUp attempts to run all up migrations against our Postgres instance.
// Migrations are loaded from a go-bindata generated module that is compiled
// into our final binary.
func MigrateUp(db *sqlx.DB, logger kitlog.Logger) error {
	logger.Log("msg", "migrating up")

	m, err := getMigrator(db.DB, logger)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to execute up migrations")
	}

	return nil
}

// MigrateDown attempts to run the requested amount of down migrations against
// the database. These migrations are loaded from a go-bindata generated module
// that compiles the sql files into our binary.
func MigrateDown(db *sqlx.DB, steps int, logger kitlog.Logger) error {
	logger.Log("msg", "migrating down", "steps", steps)

	m, err := getMigrator(db.DB, logger)
	if err != nil {
		return err
	}

	err = m.Steps(-steps)
	if err != nil {
		return errors.Wrap(err, "failed to execute down migrations")
	}

	return nil
}

// MigrateDownAll attempts to run all down migrations to roll the database back
// to its initial state.
func MigrateDownAll(db *sqlx.DB, logger kitlog.Logger) error {
	m, err := getMigrator(db.DB, logger)
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil {
		return errors.Wrap(err, "failed to execute all down migrations")
	}

	return nil
}

// NewMigration is a function that is intended to create a pair of correctly
// named migration files within the nominated migration directory. The created
// files will be empty.
func NewMigration(dirName, migrationName string, logger kitlog.Logger) error {
	if migrationName == "" {
		return errors.New("Must specify a name when creating a migration")
	}

	re := regexp.MustCompile(`\A[a-zA-Z]+\z`)
	if !re.MatchString(migrationName) {
		return errors.New("Name must be a single CamelCased string with no numbers or special characters")
	}

	migrationID := time.Now().Format("20060102150405") + "_" + snaker.CamelToSnake(migrationName)
	upFileName := fmt.Sprintf("%s.up.sql", migrationID)
	downFileName := fmt.Sprintf("%s.down.sql", migrationID)

	logger.Log("upFile", upFileName, "downFile", downFileName, "directory", dirName, "msg", "creating migration files")

	err := makeFile(dirName, upFileName)
	if err != nil {
		return err
	}

	err = makeFile(dirName, downFileName)
	if err != nil {
		return err
	}

	return nil
}

// makeFile is a helper function that creates a file within the specified
// directory. This is an empty file into which the user should then write their
// sql migration.
func makeFile(dirName, fileName string) error {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to make directory for migrations")
	}

	path := filepath.Join(dirName, fileName)

	f, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "failed to make file: %s", path)
	}

	return f.Close()
}

// getMigrator returns a new migrate.Migrate instance or an error. If we have
// managed to create our migrator it will be configured to read migrations from
// bindata assets packaged with the application.
func getMigrator(db *sql.DB, logger kitlog.Logger) (*migrate.Migrate, error) {
	dbDriver, err := psm.WithInstance(db, &psm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create postgres migrator")
	}

	source := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		},
	)

	sourceDriver, err := bindata.WithInstance(source)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create source driver")
	}

	migrator, err := migrate.NewWithInstance(
		"go-bindata",
		sourceDriver,
		"postgres",
		dbDriver,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create migrator")
	}

	migrator.Log = newLogAdapter(logger, true)

	return migrator, nil
}

// newLogAdapter is a helper function that wraps the go-kit logger into our
// custom logAdapter type which allows it to be used by the migration package.
func newLogAdapter(logger kitlog.Logger, verbose bool) migrate.Logger {
	return &logAdapter{
		logger:  logger,
		verbose: verbose,
	}
}

// logAdapter is our simple wrapper type around the go-kit logger to make it
// adhere to the interface the migration package requires.
type logAdapter struct {
	logger  kitlog.Logger
	verbose bool
}

// Printf is the function that takes the same input as fmt.Printf, and bundles
// that output into a message the go-kit Logger outputs.
func (l *logAdapter) Printf(format string, v ...interface{}) {
	l.logger.Log("msg", fmt.Sprintf(format, v...))
}

// Verbose retuns true when verbose logging output is required.
func (l *logAdapter) Verbose() bool {
	return l.verbose
}
