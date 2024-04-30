package main

import (
	"context"
	"flag"
	"fmt"
	pb "hml/protos/gen/protos"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type leaseServer struct {
	pb.UnimplementedLeaseServiceServer
	mu sync.Mutex
}

func newServer() *leaseServer {
	s := &leaseServer{}
	return s
}

func (s *leaseServer) CreateLease(ctx context.Context, request *pb.CreateLeaseRequest) (*pb.CreateLeaseResponse, error) {
	return &pb.CreateLeaseResponse{ClientId: request.ClientId, Key: request.Key, Namespace: request.Namespace}, nil
}

var (
	myAddr = flag.String("address", "localhost:50051", "TCP host+port for this node")
	raftID = flag.String("raft_id", "", "Node ID used by Raft")

	raftDir       = flag.String("raft_data_dir", "data/", "Raft data dir")
	raftBootstrap = flag.Bool("raft_bootstrap", false, "Whether to bootstrap the Raft cluster")
)

func main() {
	flag.Parse()

	if *raftID == "" {
		log.Fatalf("flag raftID is required")
	}

	//ctx := context.Background()
	_, port, err := net.SplitHostPort(*myAddr)
	if err != nil {
		log.Fatalf("failed to parse local address (%q): %v", *myAddr, err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err != nil {
		log.Fatalf("failed to start raft: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterLeaseServiceServer(s, &leaseServer{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
