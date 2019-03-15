package rpc_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/thingful/simular"
	twirp "github.com/thingful/twirp-policystore-go"

	"github.com/DECODEproject/iotpolicystore/pkg/config"
	"github.com/DECODEproject/iotpolicystore/pkg/postgres"
	"github.com/DECODEproject/iotpolicystore/pkg/rpc"
)

type PolicyStoreSuite struct {
	suite.Suite
	ps twirp.PolicyStore
}

// component is a little interface to allow us to call Start()/Stop() on the rpc
// instance which we otherwise can treat as an implementation of the twirp
// generated interface
type component interface {
	Start() error
	Stop() error
}

func (s *PolicyStoreSuite) SetupTest() {
	logger := kitlog.NewNopLogger()
	connStr := os.Getenv("POLICYSTORE_DATABASE_URL")

	config := &config.Config{
		ConnStr:            connStr,
		EncryptionPassword: "password",
		Logger:             logger,
		ClientTimeout:      1,
		DashboardURL:       "http://bcnnow.decodeproject.eu",
	}

	db := postgres.NewDB(config)
	err := db.Start()
	if err != nil {
		s.T().Fatalf("failed to start db: %v", err)
	}

	postgres.MigrateDownAll(db.DB, logger)
	postgres.MigrateUp(db.DB, logger)

	s.ps = rpc.NewPolicyStore(config, db)

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
	simular.Activate()
	defer simular.DeactivateAndReset()

	simular.RegisterStubRequests(
		simular.NewStubRequest(
			http.MethodPost,
			"http://bcnnow.decodeproject.eu/community/create_encrypted",
			simular.NewStringResponder(200, `{"id":12,"public_key": "community_key"}`),
		),
	)

	req := &twirp.CreateEntitlementPolicyRequest{
		Label: "policy label",
		Operations: []*twirp.Operation{
			&twirp.Operation{
				SensorId: 2,
				Action:   twirp.Operation_SHARE,
			},
		},
		AuthorizableAttributeId:     "abc123",
		CredentialIssuerEndpointUrl: "http://credential.com",
	}

	assert.NotNil(s.T(), req)

	createResp, err := s.ps.CreateEntitlementPolicy(context.Background(), req)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), "", createResp.CommunityId)
	assert.NotEqual(s.T(), "", createResp.Token)

	listResp, err := s.ps.ListEntitlementPolicies(context.Background(), &twirp.ListEntitlementPoliciesRequest{})
	assert.Nil(s.T(), err)
	assert.Len(s.T(), listResp.Policies, 1)

	policy := listResp.Policies[0]
	assert.Equal(s.T(), createResp.CommunityId, policy.CommunityId)
	assert.Len(s.T(), policy.Operations, 1)
	assert.Equal(s.T(), "abc123", policy.AuthorizableAttributeId)
	assert.Equal(s.T(), "http://credential.com", policy.CredentialIssuerEndpointUrl)

	operation := policy.Operations[0]
	assert.Equal(s.T(), twirp.Operation_SHARE, operation.Action)
	assert.Equal(s.T(), uint32(2), operation.SensorId)

	_, err = s.ps.DeleteEntitlementPolicy(context.Background(), &twirp.DeleteEntitlementPolicyRequest{
		CommunityId: createResp.CommunityId,
		Token:       createResp.Token,
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
			label: "missing credential issuer url",
			request: &twirp.CreateEntitlementPolicyRequest{
				Label:                       "foo",
				AuthorizableAttributeId:     "abc123",
				CredentialIssuerEndpointUrl: "",
			},
			expectedError: "twirp error invalid_argument: credential_issuer_endpoint_url is required",
		},
		{
			label: "missing label",
			request: &twirp.CreateEntitlementPolicyRequest{
				Label:                       "",
				AuthorizableAttributeId:     "abc123",
				CredentialIssuerEndpointUrl: "http://credential.com",
			},
			expectedError: "twirp error invalid_argument: label is required",
		},
		{
			label: "missing authorizable attribute id",
			request: &twirp.CreateEntitlementPolicyRequest{
				Label:                       "foo",
				AuthorizableAttributeId:     "",
				CredentialIssuerEndpointUrl: "http://credential.com",
			},
			expectedError: "twirp error invalid_argument: authorizable_attribute_id is required",
		},
		{
			label: "bins for share",
			request: &twirp.CreateEntitlementPolicyRequest{
				Label:                       "foobar",
				AuthorizableAttributeId:     "abc123",
				CredentialIssuerEndpointUrl: "http://credential.com",
				Operations: []*twirp.Operation{
					&twirp.Operation{
						SensorId: 2,
						Action:   twirp.Operation_SHARE,
						Bins:     []float64{0, 10},
					},
				},
			},
			expectedError: "twirp error invalid_argument: operation SHARE type must not specify bins or an interval",
		},
		{
			label: "interval for share",
			request: &twirp.CreateEntitlementPolicyRequest{
				Label:                       "foobar",
				AuthorizableAttributeId:     "abc123",
				CredentialIssuerEndpointUrl: "http://credential.com",
				Operations: []*twirp.Operation{
					&twirp.Operation{
						SensorId: 2,
						Action:   twirp.Operation_SHARE,
						Interval: 300,
					},
				},
			},
			expectedError: "twirp error invalid_argument: operation SHARE type must not specify bins or an interval",
		},
		{
			label: "no bins for bin",
			request: &twirp.CreateEntitlementPolicyRequest{
				Label:                       "foobar",
				AuthorizableAttributeId:     "abc123",
				CredentialIssuerEndpointUrl: "http://credential.com",
				Operations: []*twirp.Operation{
					&twirp.Operation{
						SensorId: 2,
						Action:   twirp.Operation_BIN,
					},
				},
			},
			expectedError: "twirp error invalid_argument: operation BIN type must specify bins, and no interval",
		},
		{
			label: "interval for bins",
			request: &twirp.CreateEntitlementPolicyRequest{
				Label:                       "foobar",
				AuthorizableAttributeId:     "abc123",
				CredentialIssuerEndpointUrl: "http://credential.com",
				Operations: []*twirp.Operation{
					&twirp.Operation{
						SensorId: 2,
						Action:   twirp.Operation_BIN,
						Bins:     []float64{10, 20},
						Interval: 900,
					},
				},
			},
			expectedError: "twirp error invalid_argument: operation BIN type must specify bins, and no interval",
		},
		{
			label: "no interval for moving_avg",
			request: &twirp.CreateEntitlementPolicyRequest{
				Label:                       "foobar",
				AuthorizableAttributeId:     "abc123",
				CredentialIssuerEndpointUrl: "http://credential.com",
				Operations: []*twirp.Operation{
					&twirp.Operation{
						SensorId: 2,
						Action:   twirp.Operation_MOVING_AVG,
					},
				},
			},
			expectedError: "twirp error invalid_argument: operation MOVING_AVG type must specify a non-zero positive interval, and no bins",
		},
		{
			label: "bins for moving_avg",
			request: &twirp.CreateEntitlementPolicyRequest{
				Label:                       "foobar",
				AuthorizableAttributeId:     "abc123",
				CredentialIssuerEndpointUrl: "http://credential.com",
				Operations: []*twirp.Operation{
					&twirp.Operation{
						SensorId: 2,
						Action:   twirp.Operation_MOVING_AVG,
						Interval: 900,
						Bins:     []float64{10, 20},
					},
				},
			},
			expectedError: "twirp error invalid_argument: operation MOVING_AVG type must specify a non-zero positive interval, and no bins",
		},
		{
			label: "unexpected action type",
			request: &twirp.CreateEntitlementPolicyRequest{
				Label:                       "foobar",
				AuthorizableAttributeId:     "abc123",
				CredentialIssuerEndpointUrl: "http://credential.com",
				Operations: []*twirp.Operation{
					&twirp.Operation{
						SensorId: 2,
						Action:   9,
					},
				},
			},
			expectedError: "twirp error invalid_argument: operation invalid operation type",
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
	id, err := uuid.NewRandom()
	if err != nil {
		s.T().Fatalf("failed to create id: %s", err)
	}

	testcases := []struct {
		label         string
		request       *twirp.DeleteEntitlementPolicyRequest
		expectedError string
	}{
		{
			label: "missing community_id",
			request: &twirp.DeleteEntitlementPolicyRequest{
				CommunityId: "",
				Token:       "foobar",
			},
			expectedError: "twirp error invalid_argument: community_id is required",
		},
		{
			label: "missing token",
			request: &twirp.DeleteEntitlementPolicyRequest{
				CommunityId: id.String(),
				Token:       "",
			},
			expectedError: "twirp error invalid_argument: token is required",
		},
		{
			label: "missing resource",
			request: &twirp.DeleteEntitlementPolicyRequest{
				CommunityId: id.String(),
				Token:       "foobar",
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
