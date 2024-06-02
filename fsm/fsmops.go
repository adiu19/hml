package fsm

import (
	"encoding/json"
	"errors"
	"fmt"
	storage "hml/storage"
	"io"
	"os"

	"github.com/hashicorp/raft"
)

// OperationType is a type alias
type OperationType uint8

// OperationType is a enum
const (
	GET OperationType = iota
	SET
	UPDATE
	DELETE
)

// OperationWrapper is payload sent by system when calling raft.Apply(cmd []byte, timeout time.Duration)
type OperationWrapper struct {
	Type    OperationType
	Payload interface{} // generic type to accept different payloads
}

// Apply needs to be implemented
func (f *LeaseHolderFSM) Apply(l *raft.Log) interface{} {
	switch l.Type {
	case raft.LogCommand:
		var operationPayload OperationWrapper
		if err := json.Unmarshal(l.Data, &operationPayload); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error marshalling store payload %s\n", err.Error())
			return &ResponseModel{
				Error: err,
				Data:  nil,
			}
		}

		opType := operationPayload.Type
		switch opType {
		case SET:
			operationPayload.Payload = &storage.CreateLeaseModel{}
			if err := json.Unmarshal(l.Data, &operationPayload); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error marshalling store payload %s\n", err.Error())
				return &ResponseModel{
					Error: err,
					Data:  nil,
				}
			}

			req := operationPayload.Payload.(*storage.CreateLeaseModel)
			err := f.DBAccessLayer.SetObject(req)
			if err != nil {
				// TODO: add logs
				return &ResponseModel{
					Error: err,
					Data:  nil,
				}
			}

			return &ResponseModel{
				Error: nil,
				Data:  nil,
			}
		case DELETE:
			operationPayload.Payload = &storage.GetLeaseModel{}
			if err := json.Unmarshal(l.Data, &operationPayload); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error marshalling store payload %s\n", err.Error())
				return nil
			}

			req := operationPayload.Payload.(*storage.GetLeaseModel)
			err := f.DBAccessLayer.DeleteObject(req)
			if err != nil {
				return &ResponseModel{
					Error: err,
					Data:  nil,
				}
			}
			return &ResponseModel{
				Error: nil,
				Data:  nil,
			}
		}
	}

	return nil
}

// Restore needs to be implemented
func (f *LeaseHolderFSM) Restore(rc io.ReadCloser) error {

	// clean slate
	f.DBAccessLayer.DeleteAll()

	decoder := json.NewDecoder(rc)
	for decoder.More() {
		var operationPayload OperationWrapper
		err := decoder.Decode(&operationPayload)
		if err != nil {
			return errors.New("unable to decode payload for raft respore")
		}

		opType := operationPayload.Type
		switch opType {
		case SET:
			operationPayload.Payload = &storage.CreateLeaseModel{}
			if err := decoder.Decode(&operationPayload); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error marshalling store payload %s\n", err.Error())
				return err
			}

			req := operationPayload.Payload.(*storage.CreateLeaseModel)
			err := f.DBAccessLayer.SetObject(req)
			return err
		case DELETE:
			operationPayload.Payload = &storage.GetLeaseModel{}
			if err := decoder.Decode(&operationPayload); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error marshalling store payload %s\n", err.Error())
				return err
			}

			req := operationPayload.Payload.(*storage.GetLeaseModel)
			err := f.DBAccessLayer.DeleteObject(req)
			return err
		}
	}
	return nil
}

// Snapshot needs to be implemented
func (f *LeaseHolderFSM) Snapshot() (raft.FSMSnapshot, error) {
	// Make sure that any future calls to f.Apply() don't change the snapshot.
	return &snapshot{}, nil
}

// Persist needs to be implemented
func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	return sink.Close()
}

func (s *snapshot) Release() {
}
