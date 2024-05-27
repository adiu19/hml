package cleaner

import (
	"context"
	"log"
	"time"

	"github.com/Jille/raftadmin"
	"github.com/Jille/raftadmin/proto"
	"github.com/hashicorp/raft"
)

// Run this method every 1 second
// TODO : panic is this method stops or the underlying goroutine stops for some reason, retry first
func Run(ctx context.Context, r *raft.Raft) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			// fetch state of current node
			// TODO: handle error
			state, _ := raftadmin.Get(r).State(ctx, &proto.StateRequest{})
			if raft.RaftState(state.GetState()) == raft.Leader {
				log.Println("leader node running bg task")
			}
		}

	}()

}
