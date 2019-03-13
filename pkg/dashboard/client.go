package dashboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/DECODEproject/iotpolicystore/pkg/config"
	"github.com/DECODEproject/iotpolicystore/pkg/version"
	"github.com/pkg/errors"
)

// NewClient creates a new dashboard client object with a sensible http.Client
// configured and ready to use.
func NewClient(c *config.Config) *Client {
	return &Client{
		Client: &http.Client{
			Timeout: time.Duration(c.ClientTimeout) * time.Second,
		},
		dashboardEndpoint: fmt.Sprintf("%s/community/create_encrypted", c.DashboardURL),
		userAgent:         fmt.Sprintf("%s/%s", version.BinaryName, version.VersionString()),
	}
}

// Client is a type that holds a configured HTTP client and knows how to send
// dashboard creation requests to the Dashboard API.
type Client struct {
	*http.Client

	dashboardEndpoint string
	userAgent         string
}

// createRequest is a struct used to marshal data out to the dashboard API.
type createRequest struct {
	CommunityName            string `json:"community_name"`
	CommunityID              string `json:"community_id"`
	AuthorizableAttributeID  string `json:"authorizable_attribute_id"`
	CredentialIssuerEndpoint string `json:"credential_issuer_endpoint_address"`
}

// createResponse is a struct used to marshal responses back from the dashboard
// API.
type createResponse struct {
	ID        string `json:"id"`
	PublicKey string `json:"public_key"`
}

// CreateDashboard attempts to make a call to the dashboard API to create a new
// dashboard. It returns the public key of the dashbaord or an error if there is
// any problem with the submitted values (i.e. the authorizable attribute id is
// incorrect, or the credential service is unavailable)
func (c *Client) CreateDashboard(communityID, communityName, authorizableAttributeID, credentialIssuerEndpoint string) (string, error) {
	requestObj := &createRequest{
		CommunityName:            communityName,
		CommunityID:              communityID,
		AuthorizableAttributeID:  authorizableAttributeID,
		CredentialIssuerEndpoint: credentialIssuerEndpoint,
	}

	b, err := json.Marshal(requestObj)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal request object to json")
	}

	req, err := http.NewRequest(http.MethodPost, c.dashboardEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return "", errors.Wrap(err, "failed to create http request object")
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to make http request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}

	var dashboard createResponse
	err = json.Unmarshal(body, &dashboard)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal response json")
	}

	return dashboard.PublicKey, nil
}
