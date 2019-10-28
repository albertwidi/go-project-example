# Kothak

Kothak is a `resource` holder. This library holds resource connection and object.

## Why Kothak

Instead of building another library or code to just connect to database and redis, Kothak provide `worker pool` on `redis` out of the box.

### Modern Database Architecture Configuration

Nowadays, if one or more `followers` usually sit behind a databaase connection pool or proxy such as `pgbouncer` and `haproxy` to communicate and enable load-balancing.

Kothak expect this kind of model, and only support one `follower` configuration

## SQL Database

Kothak has similar capabilities to [sqlt](https://github.com/albertwidi/sqlt) and [nap](https://github.com/tsenart/nap) by using `kothak.LoadBalancer`.

It means, all `select`, `get` and `query` will come to `follower` database. All exec is coming to `master`.

### Single Master - Follower model

Kothak only recognize one `master` database. This means Kothak is not suitable for `multi-master` database model.

### Forcing Query to Master

While [sqlt](https://github.com/albertwidi/sqlt) has `SelectMaster` function, `kothak` only providing `Master` and `Follower` function. For example:

```go
db := kothak.GetSQLDB("kosan")
db.Master().Select()
```

### Prepared Query

Kothak doesn't provide `prepare` function, and expect user to use `kothak.Master()` and `kothak.Follower()` function first to get the database. Then user can use the `preprare` function from the database object.

## Redis

To be added

### Redis Pool

To be added