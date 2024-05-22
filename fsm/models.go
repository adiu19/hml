package fsm

import (
	storage "hml/storage"
)

// LeaseHolderFSM is the fsm wrapper on top of a persistent key-value store
type LeaseHolderFSM struct {
	DB *storage.DB
}

type snapshot struct{}

// ResponseModel is the response from raft Apply
type ResponseModel struct {
	Error error
	Data  interface{}
}
