package fsm

import "hml/storage"

// LeaseHolderFSM is the fsm wrapper on top of a persistent key-value store
type LeaseHolderFSM struct {
	DBAccessLayer *storage.DBAccessLayer
}

type snapshot struct{}

// ResponseModel is the response from raft Apply
type ResponseModel struct {
	Error error
	Data  interface{}
}
