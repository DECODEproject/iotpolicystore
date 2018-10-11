package postgres_test

import (
	"os"

	kitlog "github.com/go-kit/kit/log"
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

//func (s *PostgresSuite) TestRoundTrip() {
//	req := &postgres.CreatePolicyRequest{
//		PublicKey: "abc123",
//		Label:     "policy label",
//		Operations: []postgres.Operation{
//			postgres.Operation{
//				SensorID: 2,
//				Action:   postgres.Share,
//			},
//		},
//	}
//
//	assert.NotNil(s.T(), req)
//	resp, err := s.db.CreatePolicy(req)
//	assert.Nil(s.T(), err)
//
//	// verify we have an id and token back
//	assert.NotEqual(s.T(), "", resp.ID)
//	assert.NotEqual(s.T(), "", resp.Token)
//
//	policies, err := s.db.ListPolicies()
//	assert.Nil(s.T(), err)
//	assert.Len(s.T(), policies, 1)
//
//	policy := policies[0]
//	assert.Equal(s.T(), resp.ID, policy.ID)
//	assert.Len(s.T(), policy.Operations, 1)
//
//	operation := policy.Operations[0]
//
//	assert.Equal(s.T(), postgres.Share, operation.Action)
//	assert.Equal(s.T(), 2, operation.SensorID)
//
//	deleteReq := &postgres.DeletePolicyRequest{
//		ID:    resp.ID,
//		Token: "not token",
//	}
//
//	err = s.db.DeletePolicy(deleteReq)
//	assert.NotNil(s.T(), err)
//
//	deleteReq = &postgres.DeletePolicyRequest{
//		ID:    resp.ID,
//		Token: resp.Token,
//	}
//
//	err = s.db.DeletePolicy(deleteReq)
//	assert.Nil(s.T(), err)
//}

//func TestPostgresSuite(t *testing.T) {
//	suite.Run(t, new(PostgresSuite))
//}
