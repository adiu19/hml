# Hold my Lease

A simple, distributed lease system based on the RAFT consensus protocol.

## Get Started

This work relies on using [raftadmin](https://github.com/Jille/raftadmin) to setup and manage a RAFT cluster.

### Compiling Protos
```make compile```

### Building
```make build```

This will create a binary `hml` in the project root.

### Starting Nodes
For demonstration, we'll start a locally running cluster made up of three nodes.

```
./hml --raft_bootstrap --raft_id=1 --address=localhost:50061
./hml --raft_id=2 --address=localhost:50062
./hml --raft_id=3 --address=localhost:50063
```

### Add Voters

```
raftadmin --leader multi:///localhost:50061,localhost:50062 add_voter 2 localhost:50063 0
raftadmin --leader multi:///localhost:50061,localhost:50062 add_voter 3 localhost:50063 0
raftadmin --leader multi:///localhost:50061,localhost:50062,localhost:50063 leadership_transfer
```


### Create a Lease

```
    rpc CreateLease(CreateLeaseRequest) returns (CreateLeaseResponse) {}
```

### Get a Lease

```
    rpc GetLease(GetLeaseRequest) returns (GetLeaseResponse) {}
```
