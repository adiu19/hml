./hml --raft_bootstrap --raft_id=1 --address=localhost:50061
./hml --raft_id=2 --address=localhost:50062
./hml --raft_id=3 --address=localhost:50063

raftadmin localhost:50061 add_voter 2 localhost:50062 0
raftadmin localhost:50061 add_voter 4 localhost:50064 0

raftadmin --leader multi:///localhost:50061,localhost:50062 add_voter 3 localhost:50063 0
raftadmin --leader multi:///localhost:50061,localhost:50062,localhost:50063 leadership_transfer



TODO:
1. Trying to create a lease, but it already exists [DONE]
2. Fix snapshot and log compaction [COMPLETE]
3. Locks [COMPLETE]
4. Error Channels
5. Custom Errors [COMPLETE]
6. GetLease should not require client id [DONE]
7. If BG fails, start it again, recover [DONE]
8. Created At Field [DONE]
9. Organize layers in code - right now, a lot of code is duplicated (like setting created at in response) [COMPLETE]
10. Organize models - which model belongs to db layer, which belongs to fsm layer [DONE]
11. Look at TODOs
12. Build this README
13. Test Recovery on Panic
14. Validate requests to Create [DONE]
15. Renew Flow