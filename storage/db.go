package storage

import (
	"encoding/json"
	"log"
)

// GetObject returns a key value pair
func (db *DB) GetObject(request GetLeaseModel) (*LeaseDBModel, error) {
	key := request.ClientID + "_" + request.Namespace + "_" + request.Key
	value, err := db.Store.Get([]byte(key))
	if err != nil {
		log.Fatalf("fetching db value for key %v errored out %v", key, err)
		return nil, err
	}

	obj := LeaseDBModel{}
	err = json.Unmarshal(value, &obj)
	if err != nil {
		log.Fatalf("unmarshalling db value for key %v errored out %v", key, err)
		return nil, err
	}
	return &obj, nil
}

// SetObject sets a key value pair
func (db *DB) SetObject(request CreateLeaseModel) error {
	key := request.ClientID + "_" + request.Namespace + "_" + request.Key
	objBytes, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("marshalling db value for key %v errored out %v", key, err)
		return err
	}

	err = db.Store.Set([]byte(key), objBytes)
	if err != nil {
		// TODO: add logs
		return err
	}

	return nil
}
