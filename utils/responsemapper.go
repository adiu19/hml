package utils

import (
	"fmt"
	pb "hml/protos/gen/protos"
	"hml/storage"

	"google.golang.org/protobuf/types/known/timestamppb"
)

/*
MapObjectToGetLeaseResponse maps a lease DB object to protobuf response
*/
func MapObjectToGetLeaseResponse(lease *storage.LeaseDBModel) *pb.GetLeaseResponse {
	return &pb.GetLeaseResponse{
		ClientId:     lease.ClientID,
		Key:          lease.Key,
		Namespace:    lease.Namespace,
		FencingToken: fmt.Sprint(lease.FencingToken),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: lease.CreatedAtEpochMillis,
			Nanos:   0,
		},
	}
}
