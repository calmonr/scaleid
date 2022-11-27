package graceful

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

type Graceful struct {
	group   *errgroup.Group
	timeout time.Duration
}

func New(ctx context.Context, graceful, timedout func() error, timeout time.Duration) *Graceful {
	group, ctx := errgroup.WithContext(ctx)

	g := &Graceful{
		group:   group,
		timeout: timeout,
	}

	group.Go(func() error {
		return g.listen(ctx, graceful, timedout)
	})

	return g
}

func (g Graceful) listen(ctx context.Context, graceful, timedout func() error) error {
	<-ctx.Done()

	done := make(chan struct{})

	g.group.Go(func() error {
		defer close(done)

		return graceful()
	})

	select {
	case <-done:
	case <-time.After(g.timeout):
		return timedout()
	}

	return nil
}

func (g Graceful) Action(f func() error) error {
	g.group.Go(f)

	if err := g.group.Wait(); err != nil {
		return fmt.Errorf("could not wait for graceful shutdown: %w", err)
	}

	return nil
}
