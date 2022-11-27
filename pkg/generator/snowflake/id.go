package snowflake

type ID uint64

func (s ID) Timestamp() int64 {
	return (int64(s) >> timestampShift) + twitterEpoch
}

func (s ID) DatacenterID() uint32 {
	return uint32(s) >> datacenterIDShift & MaxDatacenterID
}

func (s ID) WorkerID() uint32 {
	return uint32(s) >> workerIDShift & MaxWorkerID
}

func (s ID) Sequence() uint32 {
	return uint32(s) & MaxSequence
}
