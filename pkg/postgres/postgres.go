package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	raven "github.com/getsentry/raven-go"
	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"github.com/speps/go-hashids"
	twirp "github.com/thingful/twirp-policystore-go"
	"golang.org/x/crypto/acme/autocert"

	"github.com/DECODEproject/iotpolicystore/pkg/config"
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

	hashidData *hashids.HashIDData
	hashid     *hashids.HashID
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

	logger.Log("msg", "creating postgres DB connection", "hashidMinLength", config.HashidLength)

	hd := hashids.NewData()
	hd.Salt = config.HashidSalt
	hd.MinLength = config.HashidLength

	return &DB{
		connStr:            config.ConnStr,
		encryptionPassword: []byte(config.EncryptionPassword),
		verbose:            config.Verbose,
		logger:             logger,
		hashidData:         hd,
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

	h, err := hashids.NewWithData(d.hashidData)
	if err != nil {
		return errors.Wrap(err, "failed to create hashid generator")
	}

	d.hashid = h

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

// CreatePolicy is the function we expose from our Postgres module that is
// responsible for persisting a policy to the database. It takes our local
// PolicyRequest type which must be created from the incoming wire request, and
// returns a response struct containing the new policies id and a token the
// caller must keep secret.
func (d *DB) CreatePolicy(req *twirp.CreateEntitlementPolicyRequest) (*twirp.CreateEntitlementPolicyResponse, error) {
	// note the use of postgres native encryption to encrypt the token.
	sql := `INSERT INTO policies
    (public_key, label, token, operations)
  VALUES (:public_key, :label, pgp_sym_encrypt(:token, :encryption_password), :operations)
	RETURNING id`

	token, err := GenerateToken(TokenLength)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random token")
	}

	b, err := json.Marshal(req.Operations)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal operations JSON")
	}

	mapArgs := map[string]interface{}{
		"public_key":          req.PublicKey,
		"label":               req.Label,
		"token":               token,
		"encryption_password": d.encryptionPassword,
		"operations":          types.JSONText(b),
	}

	sql, args, err := d.DB.BindNamed(sql, mapArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to bind named query")
	}

	// note we use a Get query here so we can read back the generated ID value
	var id int
	err = d.DB.Get(&id, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute insert query")
	}

	encodedID, err := d.hashid.Encode([]int{id})
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash policy id")
	}

	return &twirp.CreateEntitlementPolicyResponse{
		PolicyId: encodedID,
		Token:    token,
	}, nil
}

// DeletePolicy is a method we expose from our Postgreo module which attempts to
// delete a policy provided the caller knows the correct id and matching token.
// It takes as input our local DeletePolicyRequest type, and if successful
// returns nothing, returning an error should any step of the operation fail.
// Note that attempting to delete a policy that has already been deleted will
// return an error to the caller.
func (d *DB) DeletePolicy(req *twirp.DeleteEntitlementPolicyRequest) error {
	// we use a CTE here to get back a count of deleted rows
	sql := `WITH deleted AS (
		DELETE FROM policies p
		WHERE p.id = :id
		AND pgp_sym_decrypt(p.token, :encryption_password) = :token
		RETURNING *)
	SELECT COUNT(*) FROM deleted`

	decodedIDList, err := d.hashid.DecodeWithError(req.PolicyId)
	if err != nil {
		return errors.Wrap(err, "failed to decode hashed id")
	}

	if len(decodedIDList) != 1 {
		return errors.New("unexpected hashed ID")
	}

	mapArgs := map[string]interface{}{
		"id":                  decodedIDList[0],
		"token":               req.Token,
		"encryption_password": d.encryptionPassword,
	}

	sql, args, err := d.DB.BindNamed(sql, mapArgs)
	if err != nil {
		return errors.Wrap(err, "failed to bind named query")
	}

	var count int
	err = d.DB.Get(&count, sql, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute delete query")
	}

	if count != 1 {
		return errors.New("no policies were deleted, either the policy id or token must be invalid")
	}

	return nil
}

// policy is an internal type used for pulling data back from the DB.
type policy struct {
	ID         int            `db:"id"`
	Label      string         `db:"label"`
	PublicKey  string         `db:"public_key"`
	Operations types.JSONText `db:"operations"`
}

// ListPolicies returns a list of all PolicyResponse structs currently
// registered in the database. We don't currently paginate or allow any
// searching or filtering of policies as it is not expected that significant
// numbers of policies will be registered.
func (d *DB) ListPolicies() ([]*twirp.ListEntitlementPoliciesResponse_Policy, error) {
	sql := `SELECT id, label, public_key, operations FROM policies ORDER BY label`

	rows, err := d.DB.Queryx(sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute read query")
	}

	policies := []*twirp.ListEntitlementPoliciesResponse_Policy{}

	for rows.Next() {
		var p policy

		err = rows.StructScan(&p)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan policy row from db")
		}

		hashedID, err := d.hashid.Encode([]int{p.ID})
		if err != nil {
			return nil, errors.Wrap(err, "failed to encode hashed id")
		}

		var operations []*twirp.Operation
		err = json.Unmarshal(p.Operations, &operations)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal operations JSON")
		}

		policyResponse := &twirp.ListEntitlementPoliciesResponse_Policy{
			PolicyId:   hashedID,
			Label:      p.Label,
			PublicKey:  p.PublicKey,
			Operations: operations,
		}

		policies = append(policies, policyResponse)
	}

	return policies, nil
}

// Ping executes the simplest query against the DB to verify that it is
// connected and responding.
func (d *DB) Ping() error {
	_, err := d.DB.Exec("SELECT 1")
	if err != nil {
		return err
	}
	return nil
}

// certificate is an internal type used for reading/writing letsencrypt
// certificates to the DB
type certificate struct {
	Key         string `db:"key"`
	Certificate []byte `db:"certificate"`
}

// Get is our implementation of the autocert.Cache interface for reading
// certificates from some underlying storage. Here we attempt to read
// certificates to Postgres.
func (d *DB) Get(ctx context.Context, key string) ([]byte, error) {
	query := `SELECT certificate FROM certificates WHERE key = $1`

	var cert []byte
	err := d.DB.Get(&cert, query, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, autocert.ErrCacheMiss
		}
		raven.CaptureError(err, map[string]string{"operation": "getCertificate"})
		return nil, errors.Wrap(err, "failed to read certificate from DB")
	}

	return cert, nil
}

// Put is our implementation of the autocert.Cache interface for writing
// certificates from underlying storage, in this case Postgres.
func (d *DB) Put(ctx context.Context, key string, data []byte) error {
	query := `INSERT INTO certificates (key, certificate)
		VALUES (:key, :certificate)
	ON CONFLICT (key)
	DO UPDATE SET certificate = EXCLUDED.certificate`

	mapArgs := map[string]interface{}{
		"key":         key,
		"certificate": data,
	}

	tx, err := d.DB.Beginx()
	if err != nil {
		raven.CaptureError(err, map[string]string{"operation": "putCertificate"})
		return errors.Wrap(err, "failed to begin transaction")
	}

	query, args, err := tx.BindNamed(query, mapArgs)
	if err != nil {
		tx.Rollback()
		raven.CaptureError(err, map[string]string{"operation": "putCertificate"})
		return errors.Wrap(err, "failed to bind named parameters")
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		raven.CaptureError(err, map[string]string{"operation": "putCertificate"})
		return errors.Wrap(err, "failed to upsert certficate")
	}

	return tx.Commit()
}

// Delete is the final method of the autocert.Cache interface that allows the
// caller to delete a certificate from the underlying store.
func (d *DB) Delete(ctx context.Context, key string) error {
	query := `DELETE FROM certificates WHERE key = $1`

	tx, err := d.DB.Beginx()
	if err != nil {
		raven.CaptureError(err, map[string]string{"operation": "deleteCertifcate"})
		return errors.Wrap(err, "failed to begin transaction")
	}

	_, err = tx.Exec(query, key)
	if err != nil {
		tx.Rollback()
		raven.CaptureError(err, map[string]string{"operation": "deleteCertificate"})
		return errors.Wrap(err, "faield to delete certificate")
	}

	return tx.Commit()
}
