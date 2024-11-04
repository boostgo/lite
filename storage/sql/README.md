# sql

### Connect (Single)
```go
const connectionString = "host=localhost port=5432 user=dbuser password=qwerty dbname=testdb sslmode=disable"
conn := sql.MustConnect(connectionString)
defer conn.Close() // or life.Tear(conn.Close)

rows, err := conn.QueryxContext(context.Background(), "SELECT * FROM users")
// ...
```

### Connect (Shard)
```go
connectionStrings := []sql.ShardConnectString{
    {
        Name:             "shard1",
        ConnectionString: "host=localhost port=5432 user=dbuser password=qwerty dbname=testdb sslmode=disable",
    },
    {
        Name:             "shard2",
        ConnectionString: "host=localhost port=5433 user=dbuser password=qwerty dbname=testdb sslmode=disable",
    },
}

connections := sql.MustConnectShards(connectionStrings, func(ctx context.Context, connections []sql.ShardConnect) sql.ShardConnect {
    // select connection by context value by key "shard"
    shard := to.String(ctx.Value("shard"))
    for _, conn := range connections {
        // if connection name equals to context "shard" value, return connection
        if conn.Name() == shard {
            return conn
        }
    }

    return nil
})
defer connections.Close() // or life.Tear(connections.Close)

client := sql.ClientShard(connections)
rows, err := client.QueryxContext(getShardContext("shard1"), "SELECT * FROM users")
// ...
```
