package snowflake_test

import (
	"testing"

	"github.com/calmonr/scaleid/pkg/generator/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	t.Parallel()

	// https://twitter.com/twitter/status/1509206476874784769
	id := snowflake.ID(1509206476874784769)

	assert.Equal(t, int64(1648657838444), id.Timestamp())
	assert.Equal(t, uint32(11), id.DatacenterID())
	assert.Equal(t, uint32(18), id.WorkerID())
	assert.Equal(t, uint32(1), id.Sequence())
}
