./hml --raft_bootstrap --raft_id=1 --address=localhost:50061
./hml --raft_id=2 --address=localhost:50062
./hml --raft_id=3 --address=localhost:50063

raftadmin localhost:50061 add_voter 2 localhost:50062 0
raftadmin --leader multi:///localhost:50061,localhost:50062 add_voter 3 localhost:50063 0
raftadmin --leader multi:///localhost:50061,localhost:50062,localhost:50063 leadership_transfer

