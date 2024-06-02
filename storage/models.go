package storage

/*
---------------------------------- Lease Profile ----------------------------------
*/

// LeaseProfile represents schema for Lease in DB
type LeaseProfile struct {
	ClientID             string
	Key                  string
	Namespace            string
	FencingToken         int64
	CreatedAtEpochMillis int64
	ExpiresAtEpochMillis int64
	UpdatedAtEpochMillis int64
}

/*
---------------------------------- Param Definitions ----------------------------------
*/

// CreateLeaseParams is the payload we receive from the raft cluster for a create request
type CreateLeaseParams struct {
	ClientID             string
	Key                  string
	Namespace            string
	ExpiresAtEpochMillis int64
	CreatedAtEpochMillis int64
}

// LeaseKeyParams is the payload we receive from the raft cluster for any operation requiring access to a lease by its key
type LeaseKeyParams struct {
	Key       string
	Namespace string
}
