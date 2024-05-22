package storage

import (
	boltdb "github.com/hashicorp/raft-boltdb"
)

// DB is a wrapper on top of a key-value store
type DB struct {
	Store boltdb.BoltStore
}

// TODO : add hierarchy here, so that base struct has client id, key, and namespace

// CreateLeaseModel is the payload we receive from the raft cluster for a create request
type CreateLeaseModel struct {
	ClientID             string
	Key                  string
	Namespace            string
	ExpiresAtEpochMillis int64
}

// GetLeaseModel is the payload we receive from the raft cluster for a get request
type GetLeaseModel struct {
	ClientID  string
	Key       string
	Namespace string
}

// LeaseDBModel represents schema for Lease in DB
type LeaseDBModel struct {
	ClientID             string
	Key                  string
	Namespace            string
	FencingToken         int64
	CreateAtEpochMillis  int64
	ExpiresAtEpochMillis int64
	UpdatedAtEpochMillis int64
}
