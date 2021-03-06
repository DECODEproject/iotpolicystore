package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"

	raven "github.com/getsentry/raven-go"
	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	twirp "github.com/thingful/twirp-policystore-go"
	"golang.org/x/crypto/acme/autocert"

	"github.com/DECODEproject/iotpolicystore/pkg/config"
	"github.com/DECODEproject/iotpolicystore/pkg/dashboard"
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

	dashboardClient *dashboard.Client
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
		dashboardClient:    dashboard.NewClient(config),
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

// CreatePolicy is the function we expose from our Postgres module that is
// responsible for persisting a policy to the database. It takes our local
// PolicyRequest type which must be created from the incoming wire request, and
// returns a response struct containing the new policies id and a token the
// caller must keep secret.
func (d *DB) CreatePolicy(req *twirp.CreateEntitlementPolicyRequest) (*twirp.CreateEntitlementPolicyResponse, error) {
	// generate the secret token we store to permit deletion of the policy
	token, err := GenerateToken(TokenLength)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random token")
	}

	// start a transaction to do the sequence of actions we need to perform in a
	// safe way
	tx, err := d.DB.Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}

	// generate uuid for the policy
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create uuid")
	}

	// call the dashboard API to create the policy and return a public key
	publicKey, err := d.dashboardClient.CreateDashboard(
		id.String(),
		req.Label,
		req.AuthorizableAttributeId,
		req.CredentialIssuerEndpointUrl,
	)

	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to create dashboard")
	}

	// now we have all the data to persist the policy
	ops, err := json.Marshal(req.Operations)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to marshal operations JSON")
	}

	descriptions, err := json.Marshal(req.Descriptions)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to marshal descriptions JSON")
	}

	// note the use of postgres native encryption to encrypt the token.
	query := `INSERT INTO policies
    (public_key, label, token, operations, authorizable_attribute_id, credential_issuer_endpoint_url, uuid, descriptions)
		VALUES (:public_key, :label, pgp_sym_encrypt(:token, :encryption_password), :operations,
		  :authorizable_attribute_id, :credential_issuer_endpoint_url, :uuid, :descriptions)`

	mapArgs := map[string]interface{}{
		"public_key":                     publicKey,
		"label":                          req.Label,
		"token":                          token,
		"encryption_password":            d.encryptionPassword,
		"operations":                     types.JSONText(ops),
		"authorizable_attribute_id":      req.AuthorizableAttributeId,
		"credential_issuer_endpoint_url": req.CredentialIssuerEndpointUrl,
		"uuid":                           id.String(),
		"descriptions":                   types.JSONText(descriptions),
	}

	query, args, err := tx.BindNamed(query, mapArgs)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to bind named query")
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to insert policy record")
	}

	return &twirp.CreateEntitlementPolicyResponse{
		CommunityId: id.String(),
		Token:       token,
	}, tx.Commit()
}

// DeletePolicy is a method we expose from our Postgreo module which attempts to
// delete a policy provided the caller knows the correct id and matching token.
// It takes as input our local DeletePolicyRequest type, and if successful
// returns nothing, returning an error should any step of the operation fail.
// Note that attempting to delete a policy that has already been deleted will
// return an error to the caller.
func (d *DB) DeletePolicy(req *twirp.DeleteEntitlementPolicyRequest) error {
	// we use a CTE here to get back a count of deleted rows
	query := `WITH deleted AS (
		DELETE FROM policies p
		WHERE p.uuid = :uuid
		AND pgp_sym_decrypt(p.token, :encryption_password) = :token
		RETURNING *)
	SELECT COUNT(*) FROM deleted`

	mapArgs := map[string]interface{}{
		"uuid":                req.CommunityId,
		"token":               req.Token,
		"encryption_password": d.encryptionPassword,
	}

	tx, err := d.DB.Beginx()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	query, args, err := tx.BindNamed(query, mapArgs)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to bind named query")
	}

	var count int
	err = tx.Get(&count, query, args...)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to execute delete query")
	}

	if count != 1 {
		tx.Rollback()
		return errors.New("no policies were deleted, either the policy id or token must be invalid")
	}

	return tx.Commit()
}

// policy is an internal type used for pulling data back from the DB.
type policy struct {
	UUID                        string             `db:"uuid"`
	Label                       string             `db:"label"`
	PublicKey                   string             `db:"public_key"`
	Operations                  types.JSONText     `db:"operations"`
	AuthorizableAttributeID     string             `db:"authorizable_attribute_id"`
	CredentialIssuerEndpointURL string             `db:"credential_issuer_endpoint_url"`
	Descriptions                types.NullJSONText `db:"descriptions"`
}

// ListPolicies returns a list of all PolicyResponse structs currently
// registered in the database. We don't currently paginate or allow any
// searching or filtering of policies as it is not expected that significant
// numbers of policies will be registered.
func (d *DB) ListPolicies() ([]*twirp.ListEntitlementPoliciesResponse_Policy, error) {
	sql := `SELECT uuid, label, public_key, operations,
		authorizable_attribute_id, credential_issuer_endpoint_url,
		descriptions
		FROM policies ORDER BY label`

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

		var operations []*twirp.Operation
		err = json.Unmarshal(p.Operations, &operations)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal operations JSON")
		}

		var descriptions map[string]string
		if p.Descriptions.Valid {
			err = json.Unmarshal(p.Descriptions.JSONText, &descriptions)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal descriptions JSON")
			}
		}

		policyResponse := &twirp.ListEntitlementPoliciesResponse_Policy{
			CommunityId:                 p.UUID,
			Label:                       p.Label,
			PublicKey:                   p.PublicKey,
			Operations:                  operations,
			AuthorizableAttributeId:     p.AuthorizableAttributeID,
			CredentialIssuerEndpointUrl: p.CredentialIssuerEndpointURL,
			Descriptions:                descriptions,
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
