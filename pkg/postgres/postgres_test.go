package postgres_test

import (
	"os"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/thingful/iotpolicystore/pkg/config"
	"github.com/thingful/iotpolicystore/pkg/postgres"
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
		s.T().Fatalf("Failed to open db connection: %v", err)
	}

	postgres.MigrateDownAll(db, logger)

	err = db.Close()
	if err != nil {
		s.T().Fatalf("Failed to close DB: %v", err)
	}

	s.db = postgres.NewDB(&config.Config{
		ConnStr:            connStr,
		HashidLength:       8,
		HashidSalt:         "hashid-salt",
		EncryptionPassword: "password",
		Logger:             logger,
	})

	err = s.db.Start()
	if err != nil {
		s.T().Fatalf("Failed to start db component: %v", err)
	}
}

func (s *PostgresSuite) TearDownTest() {
	err := s.db.Stop()
	if err != nil {
		s.T().Fatalf("Failed to stop db component: %v", err)
	}
}

func (s *PostgresSuite) TestCreatePolicy() {
	req := &postgres.PolicyRequest{
		PublicKey: "abc123",
		Label:     "policy label",
	}

	assert.NotNil(s.T(), req)
	resp, err := s.db.CreatePolicy(req)
	assert.Nil(s.T(), err)

	// verify we have an id and token back
	assert.NotEqual(s.T(), int64(0), resp.ID)
	assert.NotEqual(s.T(), "", resp.Token)
}

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}
