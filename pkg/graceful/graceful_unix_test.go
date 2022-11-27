//go:build unix

package graceful_test

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/calmonr/scaleid/internal/unittest"
	"github.com/calmonr/scaleid/pkg/graceful"
	"github.com/stretchr/testify/assert"
)

func TestGracefulAction(t *testing.T) {
	t.Parallel()

	t.Run("could not wait for graceful shutdown", func(t *testing.T) {
		t.Parallel()

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

		t.Cleanup(func() {
			stop()
		})

		err := syscall.Kill(os.Getpid(), syscall.SIGTERM)
		assert.NoError(t, err)

		g := graceful.New(
			ctx,
			func() error { return nil },
			func() error { return nil },
			1*time.Microsecond,
		)

		{
			err := g.Action(func() error {
				return unittest.ErrDummy
			})
			assert.ErrorIs(t, err, unittest.ErrDummy)
		}
	})

	t.Run("graceful", func(t *testing.T) {
		t.Parallel()

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

		t.Cleanup(func() {
			stop()
		})

		err := syscall.Kill(os.Getpid(), syscall.SIGTERM)
		assert.NoError(t, err)

		done := make(chan string)

		g := graceful.New(
			ctx,
			func() error {
				done <- "graceful"

				return nil
			},
			func() error { return nil },
			1*time.Second,
		)

		{
			err := g.Action(func() error {
				path := <-done

				assert.Equal(t, "graceful", path)

				return nil
			})
			assert.NoError(t, err)
		}
	})

	t.Run("success: timedout", func(t *testing.T) {
		t.Parallel()

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

		t.Cleanup(func() {
			stop()
		})

		err := syscall.Kill(os.Getpid(), syscall.SIGTERM)
		assert.NoError(t, err)

		done := make(chan string)
		quit := make(chan struct{})

		g := graceful.New(
			ctx,
			func() error {
				<-quit

				return nil
			},
			func() error {
				quit <- struct{}{}
				done <- "timedout"

				return nil
			},
			1*time.Microsecond,
		)

		{
			err := g.Action(func() error {
				path := <-done

				assert.Equal(t, "timedout", path)

				return nil
			})
			assert.NoError(t, err)
		}
	})
}
