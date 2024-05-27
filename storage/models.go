package storage

import "github.com/boltdb/bolt"

// DBAccessLayer is a type alias to allow definitions
type DBAccessLayer struct {
	DB *bolt.DB
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
	FencingToken int64
	ClientID     string
	Key          string
	Namespace    string
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
