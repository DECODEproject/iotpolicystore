package rpc_test

import (
	"context"
	"os"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	twirp "github.com/thingful/twirp-policystore-go"

	"github.com/thingful/iotpolicystore/pkg/config"
	"github.com/thingful/iotpolicystore/pkg/postgres"
	"github.com/thingful/iotpolicystore/pkg/rpc"
)

type PolicyStoreSuite struct {
	suite.Suite
	ps twirp.PolicyStore
}

type component interface {
	Start() error
	Stop() error
}

func (s *PolicyStoreSuite) SetupTest() {
	logger := kitlog.NewNopLogger()
	connStr := os.Getenv("POLICYSTORE_DATABASE_URL")

	db, err := postgres.Open(connStr)
	if err != nil {
		s.T().Fatalf("failed to open db connection: %v", err)
	}

	postgres.MigrateDownAll(db, logger)

	err = db.Close()
	if err != nil {
		s.T().Fatalf("failed to close db connection: %v", err)
	}

	s.ps = rpc.NewPolicyStore(&config.Config{
		ConnStr:            connStr,
		HashidLength:       8,
		HashidSalt:         "hashid-salt",
		EncryptionPassword: "password",
		Logger:             logger,
	})

	err = s.ps.(component).Start()
	if err != nil {
		s.T().Fatalf("failed to start db component: %v", err)
	}
}

func (s *PolicyStoreSuite) TearDownTest() {
	err := s.ps.(component).Stop()
	if err != nil {
		s.T().Fatalf("failed to stop db component: %v", err)
	}
}

func (s *PolicyStoreSuite) TestRoundTrip() {
	req := &twirp.CreateEntitlementPolicyRequest{
		PublicKey: "abc123",
		Label:     "policy label",
		Operations: []*twirp.Operation{
			&twirp.Operation{
				SensorId: 2,
				Action:   twirp.Operation_SHARE,
			},
		},
	}

	createResp, err := s.ps.CreateEntitlementPolicy(context.Background(), req)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), "", createResp.PolicyId)
	assert.NotEqual(s.T(), "", createResp.Token)

	listResp, err := s.ps.ListEntitlementPolicies(context.Background(), &twirp.ListEntitlementPoliciesRequest{})
	assert.Nil(s.T(), err)
	assert.Len(s.T(), listResp.Policies, 1)

	policy := listResp.Policies[0]
	assert.Equal(s.T(), createResp.PolicyId, policy.PolicyId)
	assert.Len(s.T(), policy.Operations, 1)

	operation := policy.Operations[0]
	assert.Equal(s.T(), twirp.Operation_SHARE, operation.Action)
	assert.Equal(s.T(), uint32(2), operation.SensorId)

	_, err = s.ps.DeleteEntitlementPolicy(context.Background(), &twirp.DeleteEntitlementPolicyRequest{
		PolicyId: createResp.PolicyId,
		Token:    createResp.Token,
	})

	assert.Nil(s.T(), err)

	listResp, err = s.ps.ListEntitlementPolicies(context.Background(), &twirp.ListEntitlementPoliciesRequest{})
	assert.Nil(s.T(), err)
	assert.Len(s.T(), listResp.Policies, 0)
}

func (s *PolicyStoreSuite) TestInvalidCreateRequests() {
	testcases := []struct {
		label         string
		request       *twirp.CreateEntitlementPolicyRequest
		expectedError string
	}{
		{
			label: "missing public_key",
			request: &twirp.CreateEntitlementPolicyRequest{
				PublicKey: "",
				Label:     "foo",
			},
			expectedError: "twirp error invalid_argument: public_key is required",
		},
		{
			label: "missing label",
			request: &twirp.CreateEntitlementPolicyRequest{
				PublicKey: "abc123",
				Label:     "",
			},
			expectedError: "twirp error invalid_argument: label is required",
		},
	}

	for _, tc := range testcases {
		s.T().Run(tc.label, func(t *testing.T) {
			_, err := s.ps.CreateEntitlementPolicy(context.Background(), tc.request)
			assert.NotNil(t, err)
			assert.Equal(t, tc.expectedError, err.Error())
		})
	}
}

func (s *PolicyStoreSuite) TestInvalidDeleteRequest() {
	testcases := []struct {
		label         string
		request       *twirp.DeleteEntitlementPolicyRequest
		expectedError string
	}{
		{
			label: "missing policy_id",
			request: &twirp.DeleteEntitlementPolicyRequest{
				PolicyId: "",
				Token:    "foobar",
			},
			expectedError: "twirp error invalid_argument: policy_id is required",
		},
		{
			label: "missing token",
			request: &twirp.DeleteEntitlementPolicyRequest{
				PolicyId: "abc123",
				Token:    "",
			},
			expectedError: "twirp error invalid_argument: token is required",
		},
		{
			label: "invalid policy_id",
			request: &twirp.DeleteEntitlementPolicyRequest{
				PolicyId: "abc123",
				Token:    "foobar",
			},
			expectedError: "twirp error internal: failed to decode hashed id: mismatch between encode and decode: abc123 start xm14aAYw re-encoded. result: [39775]",
		},
		{
			label: "invalid policy_id (double hashid)",
			request: &twirp.DeleteEntitlementPolicyRequest{
				PolicyId: "Vbg3HEbX",
				Token:    "foobar",
			},
			expectedError: "twirp error internal: unexpected hashed ID",
		},
		{
			label: "missing resource",
			request: &twirp.DeleteEntitlementPolicyRequest{
				PolicyId: "xm14aAYw",
				Token:    "foobar",
			},
			expectedError: "twirp error internal: no policies were deleted, either the policy id or token must be invalid",
		},
	}

	for _, tc := range testcases {
		s.T().Run(tc.label, func(t *testing.T) {
			_, err := s.ps.DeleteEntitlementPolicy(context.Background(), tc.request)
			assert.NotNil(t, err)
			assert.Equal(t, tc.expectedError, err.Error())
		})
	}
}

func TestPolicyStoreSuite(t *testing.T) {
	suite.Run(t, new(PolicyStoreSuite))
}
