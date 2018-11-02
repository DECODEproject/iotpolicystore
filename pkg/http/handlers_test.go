package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DECODEproject/iotpolicystore/pkg/config"
	"github.com/stretchr/testify/assert"

	kitlog "github.com/go-kit/kit/log"

	"github.com/DECODEproject/iotpolicystore/pkg/postgres"
)

func TestHealthCheckHandler(t *testing.T) {
	connStr := os.Getenv("POLICYSTORE_DATABASE_URL")
	logger := kitlog.NewNopLogger()

	db := postgres.NewDB(&config.Config{
		ConnStr:            connStr,
		Logger:             logger,
		HashidLength:       8,
		HashidSalt:         "salt",
		EncryptionPassword: "password",
	})

	err := db.Start()
	assert.Nil(t, err)

	defer db.Stop()

	req, err := http.NewRequest(http.MethodGet, "/pulse", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := healthCheckHandler(db)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Unexpected response code, got %v, expected %v", status, http.StatusOK)
	}

	if rr.Body.String() != "ok" {
		t.Errorf("Unexpected response body, got %s, expected %s", rr.Body.String(), "ok")
	}
}
