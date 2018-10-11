package rpc

import (
	"context"

	"github.com/twitchtv/twirp"

	kitlog "github.com/go-kit/kit/log"
	ps "github.com/thingful/twirp-policystore-go"

	"github.com/thingful/iotpolicystore/pkg/config"
	"github.com/thingful/iotpolicystore/pkg/postgres"
)

// policystore is the type that we use to implement the PolicyStore interface
type policystore struct {
	logger  kitlog.Logger
	verbose bool
	db      *postgres.DB
}

// ensure we conform to the interface at compile time
var _ ps.PolicyStore = &policystore{}

// NewPolicyStore returns a new policy store instance. It is not ready to be
// used until the Start() method is called on the object.
func NewPolicyStore(config *config.Config) ps.PolicyStore {
	db := postgres.NewDB(config)

	logger := kitlog.With(config.Logger, "module", "rpc")
	logger.Log("msg", "creating policystore rpc server")

	return &policystore{
		logger:  logger,
		verbose: config.Verbose,
		db:      db,
	}
}

// Start starts the policystore and any child components running
func (p *policystore) Start() error {
	p.logger.Log("msg", "starting policystore rpc server")

	return p.db.Start()
}

// Stop stops the policystore and any child components from running
func (p *policystore) Stop() error {
	p.logger.Log("msg", "stopping policystore rpc server")

	return p.db.Stop()
}

// CreateEntitlementPolicy is our implementation of one of the methods defined
// in our Twirp interface. This method is the mechanism by which a caller can
// write a new entitlement policy into the policystore.
func (p *policystore) CreateEntitlementPolicy(ctx context.Context, req *ps.CreateEntitlementPolicyRequest) (*ps.CreateEntitlementPolicyResponse, error) {
	// validate request
	if req.PublicKey == "" {
		return nil, twirp.RequiredArgumentError("public_key")
	}

	if req.Label == "" {
		return nil, twirp.RequiredArgumentError("label")
	}

	resp, err := p.db.CreatePolicy(req)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	return resp, nil
}

// DeleteEntitlementPolicy is our implementation of one of the methods defined
// in our Twirp interface. This method is the mechanism by which a caller can
// delete a previously created entitlement policy from the datastore.
func (p *policystore) DeleteEntitlementPolicy(ctx context.Context, req *ps.DeleteEntitlementPolicyRequest) (*ps.DeleteEntitlementPolicyResponse, error) {
	if req.PolicyId == "" {
		return nil, twirp.RequiredArgumentError("policy_id")
	}

	if req.Token == "" {
		return nil, twirp.RequiredArgumentError("token")
	}

	err := p.db.DeletePolicy(req)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	return &ps.DeleteEntitlementPolicyResponse{}, nil
}

// ListEntitlementPolicies is our implementation of one of the methods defined
// in our Twirp interface. This method is the mechanism by which a caller can
// obtain a list of all registered policies suitable for presenting to an end
// user via some sort of UI component.
func (p *policystore) ListEntitlementPolicies(ctx context.Context, req *ps.ListEntitlementPoliciesRequest) (*ps.ListEntitlementPoliciesResponse, error) {
	policies, err := p.db.ListPolicies()
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	return &ps.ListEntitlementPoliciesResponse{
		Policies: policies,
	}, nil
}
