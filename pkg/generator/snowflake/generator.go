package snowflake

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/calmonr/scaleid/pkg/clock"
)

var (
	ErrClockMovingBackwards   = errors.New("the clock is moving backwards")
	ErrDatacenterIDOutOfRange = fmt.Errorf("datacenter id is out of range [0, %d]", MaxDatacenterID)
	ErrWorkerIDOutOfRange     = fmt.Errorf("worker id is out of range [0, %d]", MaxWorkerID)
)

const (
	twitterEpoch     = 1288834974657
	datacenterIDBits = 5
	workerIDBits     = 5
	sequenceBits     = 12
)

const (
	timestampShift    = sequenceBits + workerIDBits + datacenterIDBits
	datacenterIDShift = sequenceBits + workerIDBits
	workerIDShift     = sequenceBits
)

const (
	MaxDatacenterID = -1 ^ (-1 << datacenterIDBits)
	MaxWorkerID     = -1 ^ (-1 << workerIDBits)
	MaxSequence     = -1 ^ (-1 << sequenceBits)
)

type Generator struct {
	clock                  clock.Clock
	datacenterID, workerID uint32
	storedTimestamp        atomic.Int64
	currentSequence        atomic.Uint32
}

func NewGenerator(clock clock.Clock, datacenterID, workerID uint32) (*Generator, error) {
	if datacenterID > MaxDatacenterID {
		return nil, ErrDatacenterIDOutOfRange
	}

	if workerID > MaxWorkerID {
		return nil, ErrWorkerIDOutOfRange
	}

	return &Generator{
		clock:        clock,
		datacenterID: datacenterID,
		workerID:     workerID,
	}, nil
}

func (g *Generator) Generate() (ID, error) {
	storedTimestamp := g.storedTimestamp.Load()
	currentSequence := g.currentSequence.Load()

	now := g.clock.Now().UnixMilli()

	if now < storedTimestamp {
		return 0, ErrClockMovingBackwards
	}

	if currentSequence > MaxSequence {
		for storedTimestamp == now {
			now = g.clock.Now().UnixMilli()
		}
	}

	if now > storedTimestamp {
		currentSequence = 0

		g.storedTimestamp.Store(now)
		g.currentSequence.Store(currentSequence)
	}

	generated := (now - twitterEpoch) << timestampShift
	generated |= int64(g.datacenterID) << datacenterIDShift
	generated |= int64(g.workerID) << workerIDShift
	generated |= int64(currentSequence)

	g.currentSequence.Add(1)

	return ID(generated), nil
}
