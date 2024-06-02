package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

// GetObject returns a key value pair
func (dal *DBAccessLayer) GetObject(request *GetLeaseModel) (*LeaseDBModel, error) {
	key := request.Namespace + "_" + request.Key
	data := []byte{}
	err := dal.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("conf"))

		if bucket != nil {
			val := bucket.Get([]byte(key))
			// TODO : make this logic better when a key is not found
			if val != nil {
				data = append(data[:], val[:]...)
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("failed to fetch lease %v %v", key, err)
		return nil, err
	}

	if len(data) > 0 {
		obj := LeaseDBModel{}
		err = json.Unmarshal(data, &obj)
		if err != nil {
			log.Printf("unmarshalling db value for key %v and value %v errored out %v", key, data, err)
			return nil, err
		}
		return &obj, nil
	} else {
		return nil, nil
	}
}

// DeleteObject deletes a key value pair
func (dal *DBAccessLayer) DeleteObject(request *GetLeaseModel) error {
	key := request.Namespace + "_" + request.Key
	err := dal.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("conf"))
		if err != nil {
			return errors.New("unable to create bucket")
		}

		return bucket.Delete([]byte(key))
	})

	if err != nil {
		log.Printf("failed to delete lease %v %v", key, err)
		return err
	}

	return nil
}

// DeleteAll deletes all key-value pairs in the DB
func (dal *DBAccessLayer) DeleteAll() error {
	err := dal.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("conf"))
		if err != nil {
			return errors.New("unable to create bucket")
		}

		if bucket != nil {
			c := bucket.Cursor()

			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				bucket.Delete([]byte(k))
			}
		}

		return nil

	})

	if err != nil {
		log.Printf("failed to delete all leases %v", err)
		return err
	}

	return nil
}

// GetAll returns a key value pair
func (dal *DBAccessLayer) GetAll() ([]*LeaseDBModel, error) {
	var data [][]byte
	err := dal.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("conf"))

		// make sure bucket exists
		if bucket != nil {
			c := bucket.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				data = append(data, v)
				fmt.Printf("key=%s, value=%s\n", k, v)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("failed to fetch data %v", err)
		return nil, err
	}

	var response []*LeaseDBModel

	for _, byteArray := range data {
		var leaseObj LeaseDBModel
		if err := json.Unmarshal(byteArray, &leaseObj); err != nil {
			fmt.Println("Error unmarshalling byte array:", err)
			return nil, err
		}

		response = append(response, &leaseObj)
	}

	return response, nil
}

// SetObject sets a key value pair
func (dal *DBAccessLayer) SetObject(request *CreateLeaseModel) error {
	key := request.Namespace + "_" + request.Key

	err := dal.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("conf"))
		if err != nil {
			return errors.New("unable to create bucket")
		}

		data, internalErr := json.Marshal(LeaseDBModel{
			ClientID:             request.ClientID,
			Namespace:            request.Namespace,
			Key:                  request.Key,
			ExpiresAtEpochMillis: request.ExpiresAtEpochMillis,
			CreatedAtEpochMillis: request.CreatedAtEpochMillis,
		})

		if internalErr != nil {
			log.Printf("marshalling db value for key %v errored out %v", key, internalErr)
			return internalErr
		}
		existing := bucket.Get([]byte(key))
		if existing != nil {
			log.Printf("hey, lease already exists %v", key)
			return errors.New("hey, lease already exists")
		}

		internalErr = bucket.Put([]byte(key), data)
		if internalErr != nil {
			log.Printf("failed to write data %v", internalErr)
			return internalErr
		}
		return nil
	})

	if err != nil {
		// TODO: add logs
		return err
	}

	return nil
}
