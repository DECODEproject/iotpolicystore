package postgres_test

import (
	"context"
	"os"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/acme/autocert"

	"github.com/DECODEproject/iotpolicystore/pkg/config"
	"github.com/DECODEproject/iotpolicystore/pkg/postgres"
)

type PostgresSuite struct {
	suite.Suite
	db *postgres.DB
}

func (s *PostgresSuite) SetupTest() {
	logger := kitlog.NewNopLogger()
	connStr := os.Getenv("POLICYSTORE_DATABASE_URL")

	db, err := postgres.Open(connStr)
	if err != nil {
		s.T().Fatalf("Failed to open DB connection: %v", err)
	}

	postgres.MigrateDownAll(db, logger)

	err = db.Close()
	if err != nil {
		s.T().Fatalf("Failed to close DB: %v", err)
	}

	s.db = postgres.NewDB(&config.Config{
		ConnStr:       connStr,
		Logger:        logger,
		ClientTimeout: 1,
		DashboardURL:  "http://bcnnow.decodeproject.eu",
	})

	err = s.db.Start()
	if err != nil {
		s.T().Fatalf("Failed to start postgres: %v", err)
	}
}

func (s *PostgresSuite) TearDownTest() {
	s.db.Stop()
}

func (s *PostgresSuite) TestCertificates() {
	ctx := context.Background()

	// nonexistent key should report cache miss
	_, err := s.db.Get(ctx, "baz")
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), autocert.ErrCacheMiss, err)

	// write a cert
	err = s.db.Put(ctx, "foo", []byte("bar"))
	assert.Nil(s.T(), err)

	// should be able to read that cert
	cert, err := s.db.Get(ctx, "foo")
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []byte("bar"), cert)

	// should be able to delete
	err = s.db.Delete(ctx, "foo")
	assert.Nil(s.T(), err)

	// should not be gettable
	_, err = s.db.Get(ctx, "foo")
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), autocert.ErrCacheMiss, err)
}

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}
