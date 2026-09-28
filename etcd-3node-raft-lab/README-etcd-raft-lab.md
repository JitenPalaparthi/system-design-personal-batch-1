# 3-Node etcd Raft Election Lab

## Start fresh
```bash
docker compose down -v
docker compose up -d
docker compose ps
```

## Watch all logs live
```bash
docker compose logs -f --timestamps
```

## Election-focused logs
```bash
docker compose logs --no-color | grep -Ei 'election|became|leader|candidate|term|vote|raft'
```

For live filtering:
```bash
docker compose logs -f --no-color | grep -Ei --line-buffered 'election|became|leader|candidate|term|vote|raft'
```

## Cluster status
```bash
docker exec etcd1 etcdctl \
  --endpoints=http://etcd1:2379,http://etcd2:2379,http://etcd3:2379 \
  endpoint status -w table
```

Important columns: ID, IS LEADER, RAFT TERM, RAFT INDEX, RAFT APPLIED INDEX.

## Member list
```bash
docker exec etcd1 etcdctl \
  --endpoints=http://etcd1:2379 \
  member list -w table
```

## Health
```bash
docker exec etcd1 etcdctl \
  --endpoints=http://etcd1:2379,http://etcd2:2379,http://etcd3:2379 \
  endpoint health
```

## Write/read data
```bash
docker exec etcd1 etcdctl --endpoints=http://etcd1:2379 put course etcd
docker exec etcd2 etcdctl --endpoints=http://etcd2:2379 get course
docker exec etcd3 etcdctl --endpoints=http://etcd3:2379 get course
```

## Observe persistence files
```bash
docker exec etcd1 sh -c 'find /etcd-data -maxdepth 4 -type f -ls'
docker exec etcd1 sh -c 'ls -lh /etcd-data/member/wal /etcd-data/member/snap'
```

Typical important paths:
- `/etcd-data/member/wal/` — Raft write-ahead log files
- `/etcd-data/member/snap/db` — persistent bbolt backend database
- `/etcd-data/member/snap/` — snapshot-related state

Do not edit these files manually.

## Leader failure experiment
First identify the leader:
```bash
docker exec etcd1 etcdctl \
  --endpoints=http://etcd1:2379,http://etcd2:2379,http://etcd3:2379 \
  endpoint status -w table
```

In terminal 1:
```bash
docker compose logs -f --no-color | grep -Ei --line-buffered 'election|became|leader|candidate|term|vote|raft'
```

In terminal 2, stop whichever member has `IS LEADER=true`. Example only:
```bash
docker stop etcd2
```

Then check the surviving nodes. If etcd2 was the leader:
```bash
docker exec etcd1 etcdctl \
  --endpoints=http://etcd1:2379,http://etcd3:2379 \
  endpoint status -w table
```

Expected conceptually:
1. Followers stop receiving leader heartbeats.
2. An election timeout expires.
3. A follower becomes candidate/pre-candidate as applicable.
4. A new election/term occurs.
5. A majority (2 of 3 configured members) elects a new leader.
6. `RAFT TERM` increases.

Restart the old member:
```bash
docker start etcd2
```
It should rejoin as a follower and catch up rather than automatically becoming leader.

## Quorum experiment
With all 3 healthy, identify the leader. Then stop TWO members, leaving only one running:
```bash
docker stop etcd2 etcd3
```

A 3-member cluster requires quorum=2, so one isolated member cannot make progress on consensus-backed writes.

Try:
```bash
docker exec etcd1 etcdctl --endpoints=http://etcd1:2379 put noquorum test
```

Restore quorum:
```bash
docker start etcd2
```
Then restore the third:
```bash
docker start etcd3
```

## Term vs index
- **Term**: logical election epoch. It normally increases when a new election occurs.
- **Raft index**: position in the replicated Raft log; normal writes advance it.
- **Applied index**: highest committed log entry applied to the state machine.

Try several writes and compare status before/after. The indexes should advance while the term can remain unchanged as long as leadership remains stable.

## Reset everything
```bash
docker compose down -v
```
This deletes the named Docker volumes and therefore this lab's persisted etcd state.
