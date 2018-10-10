package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thingful/iotpolicystore/pkg/postgres"
)

func TestAction(t *testing.T) {
	assert.Equal(t, "Share", postgres.Share.String())
	assert.Equal(t, "Bin", postgres.Bin.String())
	assert.Equal(t, "Moving Average", postgres.Average.String())
}

func TestUnknownAction(t *testing.T) {
	assert.Equal(t, "Unknown", postgres.Action(0).String())
}
