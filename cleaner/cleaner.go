package cleaner

import (
	"context"
	"encoding/json"
	fsm "hml/fsm"
	storage "hml/storage"
	"log"
	"time"

	"github.com/Jille/raftadmin"
	"github.com/Jille/raftadmin/proto"
	"github.com/hashicorp/raft"
)

// Run this method every 1 second
func Run(ctx context.Context, r *raft.Raft, f *fsm.LeaseHolderFSM) {
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("cleaner panicked with %v, recovery initiated", r)
			}
		}()

		for range ticker.C {
			now := time.Now().UnixMilli()
			// fetch state of current node
			state, _ := raftadmin.Get(r).State(ctx, &proto.StateRequest{})
			if raft.RaftState(state.GetState()) == raft.Leader {
				log.Printf("leader running bg task")

				leases, err := f.DBAccessLayer.GetAll()
				if err != nil {
					log.Printf("unable to fetch leases %v", err)
				}

				for _, lease := range leases {
					log.Printf("evaluating lease %v %v %v %v %v ", lease.ClientID, lease.Namespace, lease.Key, lease.ExpiresAtEpochMillis, time.Now().UnixMilli())
					if lease.ExpiresAtEpochMillis*1000 <= now {
						// delete lease from raft
						log.Printf("\ndeleting lease ----------\n\tclient_id: %v \n\tnamespace: %v \n\tkey: %v \n\texpires_at: %v \n\tcurrent_time: %v", lease.ClientID, lease.Namespace, lease.Key, lease.ExpiresAtEpochMillis*1000, now)

						payload := fsm.OperationWrapper{
							Type: fsm.DELETE,
							Payload: storage.LeaseKeyParams{
								Key:       lease.Key,
								Namespace: lease.Namespace,
							},
						}

						data, err := json.Marshal(payload)
						if err != nil {
							log.Printf("unable to expire lease due to marshalling error %v", err)
							continue
						}

						applyFuture := r.Apply(data, 500*time.Millisecond)
						if err := applyFuture.Error(); err != nil {
							// TODO: make this better
							log.Printf("unable to expire lease due to raft apply error %v", err)
							continue
						}
					}
				}
			}
		}
	}()

}
