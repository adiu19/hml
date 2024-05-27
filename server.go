package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	fsm "hml/fsm"
	pb "hml/protos/gen/protos"
	"hml/storage"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Jille/raft-grpc-leader-rpc/leaderhealth"
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/Jille/raftadmin"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type leaseServer struct {
	pb.UnimplementedLeaseServiceServer
	mu   sync.Mutex
	raft *raft.Raft
	fsm  *fsm.LeaseHolderFSM
}

func newServer() *leaseServer {
	s := &leaseServer{}
	return s
}

func (s *leaseServer) CreateLease(ctx context.Context, request *pb.CreateLeaseRequest) (*pb.CreateLeaseResponse, error) {
	// TODO: apply log here s.raft.Apply()
	payload := fsm.OperationWrapper{
		Type: fsm.SET,
		Payload: storage.CreateLeaseModel{
			ClientID:             request.ClientId,
			Key:                  request.Key,
			Namespace:            request.Namespace,
			ExpiresAtEpochMillis: request.ExpiresAt.Seconds,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		// TODO: make this better
		return &pb.CreateLeaseResponse{}, err
	}

	applyFuture := s.raft.Apply(data, 500*time.Millisecond)
	if err := applyFuture.Error(); err != nil {
		// TODO: make this better
		return &pb.CreateLeaseResponse{}, err
	}
	return &pb.CreateLeaseResponse{ClientId: request.ClientId, Key: request.Key, Namespace: request.Namespace}, nil
}

func (s *leaseServer) GetLease(ctx context.Context, request *pb.GetLeaseRequest) (*pb.GetLeaseResponse, error) {
	xyz, err := s.fsm.DB.GetObject(&storage.GetLeaseModel{ClientID: request.ClientId, Key: request.Key, Namespace: request.Namespace})
	if err != nil {
		return &pb.GetLeaseResponse{}, err
	}

	return &pb.GetLeaseResponse{ClientId: xyz.ClientID, Key: xyz.Key, Namespace: xyz.Namespace}, nil
}

var (
	myAddr = flag.String("address", "localhost:50051", "TCP host+port for this node")

	raftID        = flag.String("raft_id", "", "Node ID used by Raft")
	raftDir       = flag.String("raft_data_dir", "data/", "Raft data dir")
	raftBootstrap = flag.Bool("raft_bootstrap", false, "Whether to bootstrap the Raft cluster")
)

func main() {
	flag.Parse()

	if *raftID == "" {
		log.Fatalf("flag raftID is required")
	}

	ctx := context.Background()
	_, port, err := net.SplitHostPort(*myAddr)
	if err != nil {
		log.Fatalf("failed to parse local address (%q): %v", *myAddr, err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// creating base directory if it doesn't exist
	createBaseDir(*raftID)

	lh := newFSM()

	r, tm, err := newRaft(ctx, *raftID, *myAddr, lh)
	if err != nil {
		log.Fatalf("failed to start raft: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterLeaseServiceServer(s, &leaseServer{raft: r, fsm: lh})

	tm.Register(s)
	leaderhealth.Setup(r, s, []string{"Example"})
	raftadmin.Register(s, r)
	reflection.Register(s)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func newFSM() *fsm.LeaseHolderFSM {
	fsmStore, err := newFSMStore(*raftID)
	if err != nil {
		log.Fatalf("failed to initialize fsm store: %v", err)
	}

	return &fsm.LeaseHolderFSM{
		DB: &storage.DB{
			Store: *fsmStore,
		},
	}
}

func createBaseDir(myID string) {
	baseDir := filepath.Join(*raftDir, myID)
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		log.Fatalf("error in creating base directory: %v", err)
	}
}

func newFSMStore(myID string) (*boltdb.BoltStore, error) {
	baseDir := filepath.Join(*raftDir, myID, "fsm")
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		log.Fatalf("error in creating fsm directory: %v", err)
	}
	fsmStore, err := boltdb.NewBoltStore(filepath.Join(baseDir, "appdata.dat"))
	if err != nil {
		return fsmStore, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "appdata.dat"), err)
	}

	return fsmStore, nil
}

func newRaft(ctx context.Context, myID, myAddress string, fsm raft.FSM) (*raft.Raft, *transport.Manager, error) {
	c := raft.DefaultConfig()
	c.LocalID = raft.ServerID(myID)

	baseDir := filepath.Join(*raftDir, myID)

	ldb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "logs.dat"))
	if err != nil {
		return nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "logs.dat"), err)
	}

	sdb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		return nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "stable.dat"), err)
	}

	fss, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		return nil, nil, fmt.Errorf(`raft.NewFileSnapshotStore(%q, ...): %v`, baseDir, err)
	}

	tm := transport.New(raft.ServerAddress(myAddress), []grpc.DialOption{grpc.WithInsecure()})

	r, err := raft.NewRaft(c, fsm, ldb, sdb, fss, tm.Transport())
	if err != nil {
		return nil, nil, fmt.Errorf("raft.NewRaft: %v", err)
	}

	if *raftBootstrap {
		cfg := raft.Configuration{
			Servers: []raft.Server{
				{
					Suffrage: raft.Voter,
					ID:       raft.ServerID(myID),
					Address:  raft.ServerAddress(myAddress),
				},
			},
		}
		f := r.BootstrapCluster(cfg)
		if err := f.Error(); err != nil {
			return nil, nil, fmt.Errorf("raft.Raft.BootstrapCluster: %v", err)
		}
	}

	return r, tm, nil
}
