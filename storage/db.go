package storage

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

// GetObject returns a key value pair
func (dal *DBAccessLayer) GetObject(request *GetLeaseModel) (*LeaseDBModel, error) {
	key := request.ClientID + "_" + request.Namespace + "_" + request.Key
	data := []byte{}
	err := dal.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("conf"))
		val := bucket.Get([]byte(key))
		log.Println(val)
		data = append(data[:], val[:]...)
		return nil
	})

	if err != nil {
		log.Printf("failed to fetch data %v", err)
		return nil, err
	}

	obj := LeaseDBModel{}
	log.Println(obj)
	err = json.Unmarshal(data, &obj)
	if err != nil {
		log.Printf("unmarshalling db value for key %v errored out %v", key, err)
		return nil, err
	}
	return &obj, nil
}

// GetAll returns a key value pair
func (dal *DBAccessLayer) GetAll() ([]*LeaseDBModel, error) {
	var data [][]byte
	err := dal.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("conf"))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			data = append(data, v)
			fmt.Printf("key=%s, value=%s\n", k, v)
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
	key := request.ClientID + "_" + request.Namespace + "_" + request.Key

	err := dal.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("conf"))
		id, _ := bucket.NextSequence()
		data, internalErr := json.Marshal(LeaseDBModel{
			ClientID:     request.ClientID,
			Namespace:    request.Namespace,
			Key:          request.Key,
			FencingToken: int64(id),
		})

		if internalErr != nil {
			log.Printf("marshalling db value for key %v errored out %v", key, internalErr)
			return internalErr
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
