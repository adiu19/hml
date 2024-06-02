./hml --raft_bootstrap --raft_id=1 --address=localhost:50061
./hml --raft_id=2 --address=localhost:50062
./hml --raft_id=3 --address=localhost:50063

raftadmin localhost:50061 add_voter 2 localhost:50062 0
raftadmin localhost:50061 add_voter 4 localhost:50064 0

raftadmin --leader multi:///localhost:50061,localhost:50062 add_voter 3 localhost:50063 0
raftadmin --leader multi:///localhost:50061,localhost:50062,localhost:50063 leadership_transfer



TODO:
1. Trying to create a lease, but it already exists [DONE]
2. Fix snapshot and log compaction
3. Locks
4. Channels
5. Custom Errors
6. GetLease should not require client id [DONE]
7. If BG fails, start it again
8. Created At Field [DONE]
9. Organize layers in code - right now, a lot of code is duplicated (like setting created at in response)
10. Organize models - which model belongs to db layer, which belongs to fsm layer
