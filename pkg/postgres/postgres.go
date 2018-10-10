package postgres

import (
	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/thingful/iotpolicystore/pkg/config"
)

const (
	// TokenLength is a constant which controls the length in bytes of the tokens
	// we generate for policies.
	TokenLength = 24
)

// DB is an exported type that exposes methods for reading/writing to our
// postgres database.
type DB struct {
	DB *sqlx.DB

	connStr            string
	encryptionPassword []byte
	verbose            bool
	logger             kitlog.Logger
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

	logger.Log("msg", "creating postgres DB connection")

	return &DB{
		connStr:            config.ConnStr,
		encryptionPassword: []byte(config.EncryptionPassword),
		verbose:            config.Verbose,
		logger:             logger,
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

	err = MigrateUp(d.DB, d.logger)
	if err != nil {
		return errors.Wrap(err, "failed to run up migrations")
	}

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

func (d *DB) CreatePolicy(policy *PolicyRequest) (*PolicyResponse, error) {
	// note the use of postgres native encryption to encrypt the token value
	sql := `INSERT INTO policies
    (public_key, label, token)
  VALUES (:public_key, :label, pgp_sym_encrypt(:token, :encryption_password))
	RETURNING id`

	token, err := GenerateToken(TokenLength)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random token")
	}

	mapArgs := map[string]interface{}{
		"public_key":          policy.PublicKey,
		"label":               policy.Label,
		"token":               token,
		"encryption_password": d.encryptionPassword,
	}

	sql, args, err := d.DB.BindNamed(sql, mapArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to bind named query")
	}

	var id int64
	err = d.DB.Get(&id, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute insert query")
	}

	return &PolicyResponse{
		ID:    id,
		Token: token,
	}, nil
}
