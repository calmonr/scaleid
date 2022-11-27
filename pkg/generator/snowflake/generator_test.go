package snowflake_test

import (
	"testing"
	"time"

	"github.com/calmonr/scaleid/pkg/clock"
	"github.com/calmonr/scaleid/pkg/generator/snowflake"
	"github.com/stretchr/testify/assert"
)

func TestNewGenerator(t *testing.T) {
	t.Parallel()

	t.Run("datacenter id is out of range", func(t *testing.T) {
		t.Parallel()

		c := clock.NewMock(time.Now())

		_, err := snowflake.NewGenerator(c, snowflake.MaxDatacenterID+1, 1)
		assert.ErrorIs(t, err, snowflake.ErrDatacenterIDOutOfRange)
	})

	t.Run("worker id is out of range", func(t *testing.T) {
		t.Parallel()

		c := clock.NewMock(time.Now())

		_, err := snowflake.NewGenerator(c, 1, snowflake.MaxWorkerID+1)
		assert.ErrorIs(t, err, snowflake.ErrWorkerIDOutOfRange)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		c := clock.NewMock(time.Now())

		_, err := snowflake.NewGenerator(c, 1, 1)
		assert.NoError(t, err)
	})
}

func TestGeneratorGenerate(t *testing.T) {
	t.Parallel()

	t.Run("the clock is moving backwards", func(t *testing.T) {
		t.Parallel()

		n := time.Now()
		c := clock.NewMock(n)

		g, err := snowflake.NewGenerator(c, 1, 1)
		assert.NoError(t, err)

		_, err = g.Generate()
		assert.NoError(t, err)

		c.Add(n.Add(-time.Millisecond))

		_, err = g.Generate()
		assert.ErrorIs(t, err, snowflake.ErrClockMovingBackwards)
	})

	t.Run("reset sequence (overflow max sequence)", func(t *testing.T) {
		t.Parallel()

		n := time.Now()
		c := clock.NewMock(n)

		g, err := snowflake.NewGenerator(c, 1, 1)
		assert.NoError(t, err)

		for i := 0; i < snowflake.MaxSequence+1; i++ {
			_, err = g.Generate()
			assert.NoError(t, err)
		}

		c.Add(n)
		c.Add(n.Add(time.Millisecond))

		id, err := g.Generate()
		assert.NoError(t, err)

		assert.Equal(t, uint32(0), id.Sequence())
	})

	t.Run("reset sequence (every millisecond)", func(t *testing.T) {
		t.Parallel()

		n := time.Now()
		c := clock.NewMock(n)

		g, err := snowflake.NewGenerator(c, 1, 1)
		assert.NoError(t, err)

		id, err := g.Generate()
		assert.NoError(t, err)

		assert.Equal(t, uint32(0), id.Sequence())

		c.Add(n.Add(time.Millisecond))

		id, err = g.Generate()
		assert.NoError(t, err)

		assert.Equal(t, uint32(0), id.Sequence())
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		n := time.Now()
		c := clock.NewMock(n)

		g, err := snowflake.NewGenerator(c, 1, 1)
		assert.NoError(t, err)

		id, err := g.Generate()
		assert.NoError(t, err)

		assert.Equal(t, n.UnixMilli(), id.Timestamp())
		assert.Equal(t, uint32(1), id.DatacenterID())
		assert.Equal(t, uint32(1), id.WorkerID())
		assert.Equal(t, uint32(0), id.Sequence())
	})
}

func BenchmarkGeneratorGenerate(b *testing.B) {
	b.Run("sequence", func(b *testing.B) {
		g, _ := snowflake.NewGenerator(clock.System{}, 1, 1)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			g.Generate() // nolint:errcheck
		}
	})

	b.Run("parallel", func(b *testing.B) {
		g, _ := snowflake.NewGenerator(clock.System{}, 1, 1)

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				g.Generate() // nolint:errcheck
			}
		})
	})
}
