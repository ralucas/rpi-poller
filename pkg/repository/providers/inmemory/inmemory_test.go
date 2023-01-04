//go:build unit

package inmemory_test

import (
	"testing"

	"github.com/ralucas/rpi-poller/internal/logging"
	"github.com/ralucas/rpi-poller/pkg/repository/providers/inmemory"
	"github.com/ralucas/rpi-poller/pkg/rpi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStockStatus(t *testing.T) {
	l := logging.NewLogger(logging.LoggerConfig{})

	store := inmemory.New(l)

	testsite := "testsite"
	testprod := "testproductname"
	teststatus := rpi.OutOfStock

	err := store.SetStockStatus(testsite, testprod, teststatus)
	require.NoError(t, err)

	status, err := store.GetStockStatus(testsite, testprod)
	require.NoError(t, err)

	assert.Equal(t, status, teststatus)
}
