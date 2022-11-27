package clock_test

import (
	"testing"
	"time"

	"github.com/calmonr/scaleid/pkg/clock"
	"github.com/stretchr/testify/assert"
)

func TestSystemNow(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		c := clock.System{}

		t1 := c.Now()

		time.Sleep(time.Nanosecond)

		t2 := c.Now()

		assert.True(t, t1.Before(t2))
	})
}

func TestMockNow(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		n := time.Now()
		c := clock.NewMock(n)

		t1 := n
		t2 := n.Add(time.Second)
		t3 := n.Add(time.Second * 2)

		c.Add(t1)
		c.Add(t2)
		c.Add(t3)

		assert.Equal(t, t1, c.Now())
		assert.Equal(t, t2, c.Now())
		assert.Equal(t, t3, c.Now())
		assert.Equal(t, t1, c.Now())
	})
}
