# storage

### Connections

One line connects
```go
sqlConn := sql.MustConnect("...")
mongoConn := mongo.Must(ctx, "", "", "localhost", 27017)
```

### Transactor

Transactor - client which helps run transactions on abstract level. <br />
In business logic layer (usecase, service) we can use transactor and not define which database is using (sql or mongo, whatever) <br />

```go
// initialize transactor
transactor := sql.NewTransactor(conn)

// wrap context with "transaction" key, then database clients will run queries by transaction
var err error
ctx, err = transactor.BeginCtx(ctx)
if err != nil {
    return err
}

// catch error and choose: rollback & commit transaction
defer func() {
    if err != nil {
        _ = transactor.RollbackCtx(ctx)
        return
    }

    _ = transactor.CommitCtx(ctx)
}()
```

### Logging

Running queries by client (sql, mongo, etc...) can log queries with their params <br />
Logging must be enabled at initializing client:
```go
func NewUserRepository(conn *sqlx.DB) *UserRepo {
	const enableLog = true
	return &UserRepo{
		conn: sql.Client(conn, enableLog),
	}
}
```

It is possible turn off logging for certain calls by providing certain "context":
```go
ctx = storage.NoLog(ctx)
users, err := userRepo.Get(ctx)
// ...
```

Also, we can check if method going to log
```go
if storage.IsNoLog(ctx) {
	fmt.Println("no storage logging")
}
```
