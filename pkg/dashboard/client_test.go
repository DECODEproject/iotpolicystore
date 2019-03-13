package dashboard_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/DECODEproject/iotpolicystore/pkg/config"
	"github.com/DECODEproject/iotpolicystore/pkg/dashboard"
	"github.com/stretchr/testify/assert"

	"github.com/thingful/simular"
)

func TestCreateOK(t *testing.T) {
	simular.Activate()
	defer simular.DeactivateAndReset()

	client := dashboard.NewClient(&config.Config{
		ClientTimeout: 1,
		DashboardURL:  "http://bcnnow.decodeproject.eu",
	})

	assert.NotNil(t, client.Client)

	simular.RegisterStubRequests(
		simular.NewStubRequest(
			http.MethodPost,
			"http://bcnnow.decodeproject.eu/community/create_encrypted",
			simular.NewStringResponder(200, `{"id":"abc123","public_key": "foobarkey"}`),
			simular.WithBody(
				bytes.NewBufferString(`{"community_name":"name","community_id":"id","authorizable_attribute_id":"attribute_id","credential_issuer_endpoint_address":"http://credential.com"}`),
			),
		),
	)

	publicKey, err := client.CreateDashboard("id", "name", "attribute_id", "http://credential.com")
	assert.Nil(t, err)
	assert.Equal(t, "foobarkey", publicKey)

	assert.Nil(t, simular.AllStubsCalled())
}
