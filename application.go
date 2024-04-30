package main

import (
	"io"

	"github.com/hashicorp/raft"
)

// LeaseHolder is temp
type LeaseHolder struct {
}

// var _ raft.FSM = &LeaseHolder{}

type snapshot struct{}

// Apply needs to be implemented
func (f *LeaseHolder) Apply(l *raft.Log) interface{} {
	return nil
}

// Restore needs to be implemented
func (f *LeaseHolder) Restore(r io.ReadCloser) error {
	return nil
}

// Snapshot needs to be implemented
func (f *LeaseHolder) Snapshot() (raft.FSMSnapshot, error) {
	// Make sure that any future calls to f.Apply() don't change the snapshot.
	return &snapshot{}, nil
}

// Persist needs to be implemented
func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	return sink.Close()
}

func (s *snapshot) Release() {
}
