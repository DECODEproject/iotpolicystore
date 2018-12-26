package rpc

import (
	"context"

	"github.com/twitchtv/twirp"

	kitlog "github.com/go-kit/kit/log"
	ps "github.com/thingful/twirp-policystore-go"

	"github.com/DECODEproject/iotpolicystore/pkg/config"
	"github.com/DECODEproject/iotpolicystore/pkg/postgres"
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
func NewPolicyStore(config *config.Config, db *postgres.DB) ps.PolicyStore {
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
	if p.verbose {
		p.logger.Log("publicKey", req.PublicKey, "label", req.Label, "msg", "createPolicy")
	}

	// validate request
	if req.PublicKey == "" {
		return nil, twirp.RequiredArgumentError("public_key")
	}

	if req.Label == "" {
		return nil, twirp.RequiredArgumentError("label")
	}

	err := validateOperations(req.Operations)
	if err != nil {
		return nil, err
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
	if p.verbose {
		p.logger.Log("policyID", req.PolicyId, "msg", "deletePolicy")
	}

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
	if p.verbose {
		p.logger.Log("msg", "listPolicies")
	}

	policies, err := p.db.ListPolicies()
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}

	return &ps.ListEntitlementPoliciesResponse{
		Policies: policies,
	}, nil
}

// validateOperations validates the content of all operations
func validateOperations(operations []*ps.Operation) error {
	for _, operation := range operations {
		err := validateOperation(operation)
		if err != nil {
			return err
		}
	}
	return nil
}

// validateOperation validates a single operation. Used by validateOperations to
// check each operation. An operation is valid if: 1) it is of type SHARE and
// has no bins or interval specified, 2) if it is of type BIN and it has bins
// specified, 3) it is of type MOVING_AVG and it has a time interval. Any other
// combination should be flagged as invalid.
func validateOperation(operation *ps.Operation) error {
	switch operation.Action {
	case ps.Operation_SHARE:
		if len(operation.Bins) > 0 || operation.Interval > 0 {
			return twirp.InvalidArgumentError("operation", "SHARE type must not specify bins or an interval")
		}
	case ps.Operation_BIN:
		if len(operation.Bins) == 0 || operation.Interval > 0 {
			return twirp.InvalidArgumentError("operation", "BIN type must specify bins, and no interval")
		}
	case ps.Operation_MOVING_AVG:
		if operation.Interval <= 0 || len(operation.Bins) > 0 {
			return twirp.InvalidArgumentError("operation", "MOVING_AVG type must specify a non-zero positive interval, and no bins")
		}
	default:
		return twirp.InvalidArgumentError("operation", "invalid operation type")
	}
	return nil
}
