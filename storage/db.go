package leasestorage

// LeaseDBObj represents the schema of leases
type LeaseDBObj struct {
	ClientID             string
	Key                  string
	Namespace            string
	CreateAtEpochMillis  int64
	ExpiresAtEpochMillis int64
}

func CreateLease() *LeaseDBObj {
	return &LeaseDBObj{}
}
